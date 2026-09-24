#!/bin/sh
# Copyright (c) 2026 Andrew David LeTourneau
# SPDX-License-Identifier: MIT OR Zlib
set -eu

source_dir=$(pwd)
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
mkdir -p "$tmp/scripts" "$tmp/releases" "$tmp/bin"
cp "$source_dir/scripts/check-version.sh" "$source_dir/scripts/create-release-tag.sh" "$tmp/scripts/"
cd "$tmp"
unset CI_COMMIT_TAG
git init -q
git config user.name 'Release tests'
git config user.email 'release-tests@example.invalid'
git -c commit.gpgsign=false commit -qm initial --allow-empty
export CI_COMMIT_SHA="$(git rev-parse HEAD)"
export CI_DEFAULT_BRANCH=main CI_COMMIT_BRANCH=main CI_PIPELINE_SOURCE=push

pass() { "$@" >test-output 2>&1 || { cat test-output; exit 1; }; }
fail() { if "$@" >test-output 2>&1; then echo "Unexpected success: $*" >&2; exit 1; fi; }

printf '0.1.0\n' >VERSION
fail sh scripts/check-version.sh
printf 'Release notes\n' >releases/v0.1.0.md
pass sh scripts/check-version.sh
pass env CI_COMMIT_TAG=v0.1.0 sh scripts/check-version.sh
fail env CI_COMMIT_TAG=v0.2.0 sh scripts/check-version.sh
for value in v0.1.0 01.1.0 0.1 0.1.0-rc.1 '0.1.0+build' ''; do
  printf '%s\n' "$value" >VERSION
  fail sh scripts/check-version.sh
done
printf '0.1.0\n0.2.0\n' >VERSION
fail sh scripts/check-version.sh
printf '0.1.0\n' >VERSION
git tag v0.2.0
fail sh scripts/check-version.sh
git tag -d v0.2.0 >/dev/null

# Fake the transport: no test can reach GitLab or create a real release tag.
cat >bin/curl <<'CURL'
#!/bin/sh
printf '%s\n' "$@" >curl-args
exit "${CURL_EXIT:-0}"
CURL
chmod +x bin/curl
export PATH="$tmp/bin:$PATH"
export RELEASE_TOKEN=test-token CI_API_V4_URL=https://example.invalid/api/v4 CI_PROJECT_ID=123
fail env CI_COMMIT_BRANCH=feature sh scripts/create-release-tag.sh
fail env CI_PIPELINE_SOURCE=merge_request_event sh scripts/create-release-tag.sh
fail env RELEASE_TOKEN= sh scripts/create-release-tag.sh
test ! -e curl-args
pass sh scripts/create-release-tag.sh
grep -Fx "tag_name=v0.1.0" curl-args >/dev/null
grep -Fx "ref=$CI_COMMIT_SHA" curl-args >/dev/null
grep -Fx 'https://example.invalid/api/v4/projects/123/repository/tags' curl-args >/dev/null
fail env CURL_EXIT=22 sh scripts/create-release-tag.sh
rm curl-args
git tag v0.1.0
pass env RELEASE_TOKEN= sh scripts/create-release-tag.sh
test ! -e curl-args

# Existing tags must belong to the pipeline commit's history.
git -c commit.gpgsign=false commit -qm unrelated --allow-empty
git tag -f v0.1.0 >/dev/null
fail sh scripts/create-release-tag.sh
test ! -e curl-args
echo 'Release version and tag tests passed.'
