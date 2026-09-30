#!/bin/sh
# Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
set -eu
# semantic-release prepares archives before creating its tag. Snapshot mode
# skips tagging/publishing; the explicit version is still embedded in binaries.
export PATCHBOARD_VERSION="${PATCHBOARD_VERSION:-0.0.0-dev.$(git rev-parse --short HEAD)}"
./.cache/bin/goreleaser check
./.cache/bin/goreleaser release --snapshot --clean
