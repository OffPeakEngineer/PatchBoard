#!/bin/sh
# Copyright (c) 2026 Andrew David LeTourneau
# SPDX-License-Identifier: MIT OR Zlib
set -eu

version=$(cat VERSION)
if ! printf '%s\n' "$version" | grep -Eq '^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$' || [ "$(printf '%s\n' "$version" | wc -l)" -ne 1 ]; then
  echo 'VERSION must contain a stable SemVer version, such as 0.1.0.' >&2
  exit 1
fi
tag="v$version"
if [ -n "${CI_COMMIT_TAG:-}" ] && [ "$CI_COMMIT_TAG" != "$tag" ]; then
  echo "Tag $CI_COMMIT_TAG does not match VERSION ($tag)." >&2
  exit 1
fi
if [ ! -s "releases/$tag.md" ]; then
  echo "Missing release notes: releases/$tag.md" >&2
  exit 1
fi
# A full fetch (GIT_DEPTH=0 in CI) makes the comparison include all releases.
latest=$(git tag --list | grep -E '^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$' | sort -V | tail -n 1)
if [ -n "$latest" ] && [ "$(printf '%s\n%s\n' "$latest" "$tag" | sort -V | tail -n 1)" != "$tag" ]; then
  echo "VERSION ($tag) is older than the latest release ($latest)." >&2
  exit 1
fi
