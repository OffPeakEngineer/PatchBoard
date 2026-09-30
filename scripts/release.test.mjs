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
