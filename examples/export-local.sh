#!/bin/sh
# Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
# Run from your project's root, or pass explicit task and output paths.
set -eu

PATCHBOARD_REF=${PATCHBOARD_REF:-097f0d865eec8c7e0151efd270132b5cf4f95a06}
PATCHBOARD_SOURCE=${PATCHBOARD_SOURCE:-https://github.com/OffPeakEngineer/PatchBoard.git}
task_directory=${1:-tasks}
output=${2:-public/index.html}
# Resolve paths before entering the isolated tooling checkout; spaces are allowed.
task_directory=$(cd "$task_directory" && pwd)
case "$output" in
  /*) ;;
  *) output="$PWD/$output" ;;
esac
if [ ! -f "$task_directory/kanban.html" ]; then
  echo "Missing $task_directory/kanban.html; initialize or update this board first." >&2
  exit 1
fi

tooling=$(mktemp -d)
trap 'rm -rf "$tooling"' EXIT
trap 'exit 1' HUP INT TERM
git -C "$tooling" init --quiet
git -C "$tooling" remote add origin "$PATCHBOARD_SOURCE"
git -C "$tooling" fetch --quiet --depth=1 origin "$PATCHBOARD_REF"
git -C "$tooling" checkout --quiet --detach FETCH_HEAD
(
  cd "$tooling"
  npm ci
  # On a Linux CI runner, install OS libraries once with:
  # npx playwright install --with-deps chromium
  npx playwright install chromium
  node scripts/export-board.mjs "$task_directory" "$output"
)
printf 'Board saved to %s\n' "$output"
