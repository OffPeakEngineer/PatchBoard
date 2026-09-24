<!-- Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib -->
# Releases

The canonical source and release host is
[GitLab](https://gitlab.com/off-peak.engineer/utilities/patchboard).
The module path is `gitlab.com/off-peak.engineer/utilities/patchboard`. The project
is currently private: source, package, and release downloads require access.

## Version and publication flow

1. Change `VERSION` to a stable SemVer version without a `v` prefix, and add
   `releases/vX.Y.Z.md` with the release notes in the same MR. The first candidate
   is `0.1.0`; later release versions must increase.
2. MR CI validates the version and notes, runs race tests, vet, board lint, and
   release-script tests, then builds all six platform archives with GoReleaser.
   Snapshot artifacts are downloadable from the pipeline for one week.
3. After merge, the default-branch push pipeline repeats those checks and creates
   an annotated `vX.Y.Z` tag at its exact commit using the GitLab Tags API. An
   unchanged version whose tag is already in the commit's history is a no-op.
4. The tag pipeline tests again, checks the tag against `VERSION`, and publishes
   a GitLab release using the matching notes. GoReleaser stores archives and
   checksums in the Generic Package Registry and attaches download links to the
   release. Those downloads do not depend on the one-week CI artifact expiry.

Merging a version bump authorizes publication. Do not enable auto-merge for a
release MR until its version, notes, and pipeline setup have been reviewed.
Ordinary main commits with an unchanged version do not produce another release.
MR, feature-branch, schedule, and manually started branch pipelines do not tag.
Only stable `vX.Y.Z` tag pipelines publish; prereleases are not yet configured.

## One-time GitLab setup

- Enable the Package Registry (already enabled in the current project).
- Protect the default branch and `v*` tags. Allow the release token's identity
  to create those protected tags.
- Add `RELEASE_TOKEN` as a masked, protected CI variable: a project access token
  with `api` scope and a role permitted to create release tags. This token is
  used only by `prepare-release`. A personal API token with equivalent access is
  an alternative if project access tokens are unavailable.
- The tag publisher uses the built-in `CI_JOB_TOKEN`, passed as `GITLAB_TOKEN`.
  No GitHub or package-manager credentials are required.

A job token cannot create tags through the Tags API, and Git pushes using a job
token do not start another pipeline. The separate API token is therefore needed
for automatic tag creation. See the [GitLab job-token permissions](https://docs.gitlab.com/ci/jobs/ci_job_token/)
and [GoReleaser GitLab configuration](https://goreleaser.com/customization/publish/scm/gitlab/).

## Artifacts and local validation

| OS | Architectures | Archive |
| --- | --- | --- |
| Linux | amd64, arm64 | `.tar.gz` |
| macOS (`darwin`) | amd64, arm64 | `.tar.gz` |
| Windows | amd64, arm64 | `.zip` containing `patchboard.exe` |

Archive names are `patchboard_X.Y.Z_OS_ARCH`; each includes the binary, README,
and all license files. `checksums.txt` lists SHA-256 hashes. On Linux use
`sha256sum --check --ignore-missing checksums.txt` after downloading archives;
on macOS use `shasum -a 256`, or on Windows `Get-FileHash -Algorithm SHA256`,
and compare the result with the matching checksum entry.

CI uses Go 1.25, matching the minimum in `go.mod`, and GoReleaser 2.18.2.
The installer verifies the downloaded tool against its published checksum.
On a Linux amd64/arm64 development machine:

```sh
sh scripts/check-version.sh
sh scripts/test-release.sh
go test -race ./...
go vet ./...
go run ./cmd/patchboard lint
sh scripts/install-goreleaser.sh
.cache/bin/goreleaser check
.cache/bin/goreleaser release --snapshot --clean
```

Snapshot builds do not publish or need release credentials. Test the native
binary's `--version`, `--help`, and board initialization from outside this
checkout. Cross-compilation verifies buildability; maintainers should also smoke
test macOS and Windows binaries and browser permission flows before publication.

## Retry and recovery

If tag creation fails because the token is missing, configure it and retry the
job. If the request reached GitLab before a network failure, confirm the existing
tag and its pipeline before retrying; never move an existing release tag.

If the tag pipeline fails, retry that pipeline's failed job after resolving the
cause. Do not delete the version tag merely to retrigger main. A tested main
commit can also be tagged manually by an authorized maintainer; its tag must
match `VERSION`. Keep published package versions out of registry cleanup rules.
Fix a published binary by releasing a new version.
