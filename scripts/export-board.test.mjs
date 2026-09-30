// Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { mkdtemp, mkdir, writeFile, copyFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { pathToFileURL } from 'node:url';
import { chromium } from 'playwright';
import { exportBoard } from './export-board.mjs';

test('exports rendered lanes, escaped task contents, and working offline links', async () => {
  const dir = await mkdtemp(join(tmpdir(), 'patchboard-export-'));
  const browser = await chromium.launch();
  try {
    await copyFile('templates/kanban.html', join(dir, 'kanban.html'));
    await writeFile(join(dir, 'board.yml'), 'states:\n  - ready\n  - done\n');
    await mkdir(join(dir, 'ready'));
    await mkdir(join(dir, 'done'));
    const markdown = '# A <script>alert(1)</script> & task\n\nFull task details <img src=x onerror=alert(1)>\n';
    await writeFile(join(dir, 'ready', 'test.md'), markdown);
    const output = join(dir, 'out', 'index.html');
    assert.equal(await exportBoard(dir, output), 1);
    const page = await browser.newPage({ javaScriptEnabled: false });
    await page.goto(pathToFileURL(output).href);
    assert.equal(await page.locator('.lane').count(), 2);
    assert.equal(await page.locator('.card').count(), 1);
    assert.equal(await page.locator('.card').getAttribute('draggable'), 'false');
    assert.equal(await page.locator('script, button, select, img').count(), 0);
    assert.equal(await page.locator('#task-0 pre').textContent(), markdown);
    await page.locator('.card a').click();
    assert.equal(new URL(page.url()).hash, '#task-0');
    // Empty lanes are a valid board; exporting an unrelated folder must fail.
    await rm(join(dir, 'ready', 'test.md'));
    assert.equal(await exportBoard(dir, output), 0);
    await writeFile(join(dir, 'board.yml'), 'states:\n  - missing\n');
    await assert.rejects(exportBoard(dir, output), /not a PatchBoard task root/);
  } finally {
    await browser.close();
    await rm(dir, { recursive: true, force: true });
  }
});
