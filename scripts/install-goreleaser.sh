#!/bin/sh
# Copyright (c) 2026 Andrew David LeTourneau
# SPDX-License-Identifier: MIT OR Zlib
set -eu

version=2.18.2
case "$(uname -m)" in
  x86_64) arch=x86_64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo 'GoReleaser installer supports Linux amd64/arm64 runners.' >&2; exit 1 ;;
esac
archive="goreleaser_Linux_$arch.tar.gz"
base="https://github.com/goreleaser/goreleaser/releases/download/v$version"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
curl --fail --silent --show-error --location "$base/$archive" -o "$tmp/$archive"
curl --fail --silent --show-error --location "$base/checksums.txt" -o "$tmp/checksums.txt"
(cd "$tmp" && grep "  $archive\$" checksums.txt | sha256sum --check --strict -)
mkdir -p .cache/bin
tar -xzf "$tmp/$archive" -C .cache/bin goreleaser
.cache/bin/goreleaser --version
