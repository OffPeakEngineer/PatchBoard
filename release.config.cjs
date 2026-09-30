// Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
module.exports = {
  branches: [process.env.RELEASE_BRANCH || 'main'],
  repositoryUrl: `https://github.com/${process.env.GITHUB_REPOSITORY || 'OffPeakEngineer/patchboard'}.git`,
  tagFormat: 'v${version}',
  plugins: [
    ['@semantic-release/commit-analyzer', { preset: 'conventionalcommits' }],
    ['@semantic-release/release-notes-generator', { preset: 'conventionalcommits' }],
    ['@semantic-release/exec', {
      prepareCmd: 'PATCHBOARD_VERSION=${nextRelease.version} sh scripts/build-archives.sh',
    }],
    ['@semantic-release/github', {
      assets: ['dist/*.tar.gz', 'dist/*.zip', 'dist/checksums.txt'],
      successComment: false,
      failComment: false,
      failTitle: false,
      releasedLabels: false,
    }],
  ],
};
