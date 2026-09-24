#!/bin/sh
# Copyright (c) 2026 Andrew David LeTourneau
# SPDX-License-Identifier: MIT OR Zlib
set -eu

: "${CI_DEFAULT_BRANCH:?Run this from a GitLab main push pipeline}"
: "${CI_COMMIT_SHA:?Missing pipeline commit}"
if [ "${CI_COMMIT_BRANCH:-}" != "$CI_DEFAULT_BRANCH" ] || [ "${CI_PIPELINE_SOURCE:-}" != push ]; then
  echo 'Release tags can only be created by a default-branch push pipeline.' >&2
  exit 1
fi
sh scripts/check-version.sh
tag="v$(cat VERSION)"
if git rev-parse --verify "refs/tags/$tag" >/dev/null 2>&1; then
  if ! git merge-base --is-ancestor "$tag" "$CI_COMMIT_SHA"; then
    echo "$tag exists outside this commit's history; refusing to reuse it." >&2
    exit 1
  fi
  echo "$tag already exists; no new release requested."
  exit 0
fi
: "${GITLAB_TOKEN:?Set a masked, protected project access token with api scope; see docs/releases.md}"
: "${CI_API_V4_URL:?Missing GitLab API URL}"
: "${CI_PROJECT_ID:?Missing GitLab project ID}"
# Use an API token: a tag pushed with CI_JOB_TOKEN does not trigger a pipeline.
# Do not retry this POST automatically; a network failure can follow creation.
curl --fail-with-body --silent --show-error --request POST \
  --header "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  --data-urlencode "tag_name=$tag" \
  --data-urlencode "ref=$CI_COMMIT_SHA" \
  --data-urlencode "message=PatchBoard $tag" \
  "$CI_API_V4_URL/projects/$CI_PROJECT_ID/repository/tags" >/dev/null
echo "Created $tag at $CI_COMMIT_SHA; its pipeline will publish the release."
