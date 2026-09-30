// Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'subject-case': [1, 'never', ['sentence-case', 'start-case', 'pascal-case', 'upper-case']],
  },
};
