// Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
import { test } from 'node:test';
import assert from 'node:assert/strict';
import { analyzeCommits } from '@semantic-release/commit-analyzer';
import config from '../release.config.cjs';

for (const [message, expected] of [
  ['fix: repair folder loading', 'patch'],
  ['feat: publish a board', 'minor'],
  ['feat!: replace the task format', 'major'],
  ['fix: migrate format\n\nBREAKING CHANGE: old tasks need migration', 'major'],
  ['docs: clarify installation', null],
  ['ci: update runners', null],
]) {
  test(`release decision for ${message.split('\n')[0]}`, async () => {
    const options = config.plugins.find(([name]) => name === '@semantic-release/commit-analyzer')[1];
    const actual = await analyzeCommits(options, {
      cwd: process.cwd(), commits: [{ hash: 'test', message }], logger: { log() {} },
    });
    assert.equal(actual, expected);
  });
}

// Exercise the writer as well as version selection: incompatible preset/writer
// majors can analyze commits successfully but fail (or emit empty release notes).
test('release notes include features, fixes, breaking changes, and GitHub links', async () => {
  const { generateNotes } = await import('@semantic-release/release-notes-generator');
  const options = config.plugins.find(([name]) => name === '@semantic-release/release-notes-generator')[1];
  const notes = await generateNotes(options, {
    cwd: process.cwd(),
    options: { repositoryUrl: config.repositoryUrl },
    commits: [
      { hash: 'a'.repeat(40), message: 'feat: export the project board' },
      { hash: 'b'.repeat(40), message: 'fix: preserve task links' },
      { hash: 'c'.repeat(40), message: 'feat!: replace the task format\n\nBREAKING CHANGE: migrate existing task metadata' },
    ],
    lastRelease: { version: '1.0.0', gitTag: 'v1.0.0' },
    nextRelease: { version: '2.0.0', gitTag: 'v2.0.0' },
    logger: { log() {} },
  });
  for (const text of ['export the project board', 'preserve task links', 'migrate existing task metadata', 'v1.0.0...v2.0.0', '/commit/']) {
    assert.ok(notes.includes(text), `release notes are missing ${text}`);
  }
});
