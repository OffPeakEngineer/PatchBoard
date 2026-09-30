// Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
import { readdir, readFile, mkdir, writeFile } from 'node:fs/promises';
import { resolve, join, dirname } from 'node:path';
import { pathToFileURL } from 'node:url';
import { chromium } from 'playwright';

// Supply read-only File System Access handles to the existing browser renderer.
// Only board configuration and Markdown files are exposed, never symlinks.
async function readTree(directory) {
  const children = [];
  for (const entry of (await readdir(directory, { withFileTypes: true })).sort((a, b) => a.name.localeCompare(b.name))) {
    if (entry.name.startsWith('.')) continue;
    if (entry.isDirectory()) {
      children.push({ name: entry.name, kind: 'directory', children: await readTree(join(directory, entry.name)) });
    } else if (entry.isFile() && (entry.name === 'board.yml' || entry.name.endsWith('.md'))) {
      children.push({ name: entry.name, kind: 'file', text: await readFile(join(directory, entry.name), 'utf8') });
    }
  }
  return children;
}

export async function exportBoard(taskDirectory = 'tasks', output = 'public/index.html') {
  const directory = resolve(taskDirectory);
  const tree = { name: 'tasks', kind: 'directory', children: await readTree(directory) };
  const browser = await chromium.launch();
  try {
    const page = await browser.newPage({ locale: 'en-US' });
    const errors = [];
    page.on('pageerror', error => errors.push(error));
    // The exporter must never restore a remembered folder or make network calls.
    await page.addInitScript(() => { delete window.showDirectoryPicker; });
    await page.route(/^https?:/, route => route.abort());
    await page.goto(pathToFileURL(join(directory, 'kanban.html')).href);
    const count = await page.evaluate(async tree => {
      function handle(node) {
        return {
          name: node.name, kind: node.kind,
          async *entries() { for (const child of node.children || []) yield [child.name, handle(child)]; },
          async getFile() { return new File([node.text], node.name); },
          async getFileHandle(name) { return child(name, 'file'); },
          async getDirectoryHandle(name) { return child(name, 'directory'); },
        };
        function child(name, kind) {
          const found = node.children?.find(value => value.name === name && value.kind === kind);
          if (!found) throw new DOMException(`Missing ${name}`, 'NotFoundError');
          return handle(found);
        }
      }
      const root = handle(tree);
      await loadHandle(root, 'readonly', false);
      document.querySelector('.toolbar').remove();
      const status = document.querySelector('#status');
      status.removeAttribute('data-i18n');
      status.textContent = 'Read-only snapshot of this project’s tasks. Updated by the project pipeline.';
      const details = document.createElement('section');
      details.id = 'task-details';
      const heading = document.createElement('h2');
      heading.textContent = 'Task files';
      details.append(heading);
      const cards = [...document.querySelectorAll('.card')];
      for (const [index, card] of cards.entries()) {
        const [state, name] = card.dataset.path.split('/');
        const file = await (await (await root.getDirectoryHandle(state)).getFileHandle(name)).getFile();
        const article = document.createElement('article');
        article.id = `task-${index}`;
        const title = document.createElement('h3');
        title.textContent = card.querySelector('.title').textContent;
        const text = document.createElement('pre');
        text.textContent = await file.text();
        const back = document.createElement('a');
        back.href = '#board';
        back.textContent = 'Back to board';
        article.append(title, text, back);
        details.append(article);
        const link = card.querySelector('a');
        link.href = `#${article.id}`;
        link.removeAttribute('target');
        card.draggable = false;
      }
      document.body.append(details);
      const style = document.createElement('style');
      style.textContent = '#task-details { padding: 24px; } #task-details article { margin-bottom: 32px; scroll-margin-top: 16px; } #task-details pre { white-space: pre-wrap; overflow-wrap: anywhere; }';
      document.head.append(style);
      // Serialization drops event listeners; removing scripts makes this a true
      // static document that also works offline with JavaScript disabled.
      document.querySelectorAll('script').forEach(script => script.remove());
      return cards.length;
    }, tree);
    if (errors.length) throw errors[0];
    const html = await page.content();
    await mkdir(dirname(resolve(output)), { recursive: true });
    await writeFile(output, html);
    return count;
  } finally {
    await browser.close();
  }
}

if (process.argv[1] && import.meta.url === pathToFileURL(resolve(process.argv[1])).href) {
  const count = await exportBoard(process.argv[2], process.argv[3]);
  console.log(`Exported ${count} task cards to ${process.argv[3] || 'public/index.html'}`);
}
