<!-- Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib -->
# Releases

[OffPeakEngineer/patchboard on GitHub](https://github.com/OffPeakEngineer/patchboard)
is the official release authority. GitLab is a backup build and release
host. Both providers test the project, build six platform archives, and publish
the rendered task board to their own Pages site. The Go module remains
`gitlab.com/off-peak.engineer/utilities/patchboard`.

## Conventional Commits and semantic-release

Use `type(scope): description` commit messages. The scope is optional:

- `fix: repair folder loading` produces a patch release.
- `feat: export a static board` produces a minor release.
- `feat!: change the task format` or a `BREAKING CHANGE:` footer produces a major release.
- `docs:`, `test:`, `chore:`, and `ci:` alone do not produce a release.

CI checks new push commits and PR/MR commit ranges with commitlint. Historical
commits are not retroactively linted. Use a Conventional Commit title when
squash-merging as well: the resulting commit determines the release.

After a tested default-branch push on GitHub, semantic-release reads commits
since the last reachable `vX.Y.Z` tag, computes the next version, generates notes,
builds archives with that exact version embedded, creates the tag, and uploads
the GitHub release assets. No npm package is published. Release jobs are
serialized. A commit with no release-worthy changes still updates Pages.

There is no manual `VERSION` bump or required release-notes file. The files in
`releases/` are historical notes. Mirror existing tags (including `v0.1.0`)
before enabling the workflow so semantic-release has the correct baseline.
Without a prior release tag, semantic-release starts at `1.0.0`.

GitLab never runs semantic-release or creates release tags. A mirrored stable
`vX.Y.Z` tag triggers its backup publisher, which builds that version and uses
`CI_JOB_TOKEN` to publish a GitLab release and Generic Package Registry assets.
Backup notes come from Git history; GitHub's semantic-release notes are official.
Do not create independent release tags on GitLab.

## Provider setup

GitHub:

1. Push the repository history and existing release tags to the official GitHub
   project. The workflow derives its repository URL and default branch from the
   running project, so no hard-coded GitHub owner is needed.
2. Enable Actions and select **GitHub Actions** as the Pages build source.
3. Permit `GITHUB_TOKEN` to write contents and create protected `v*` tags under
   your repository rules. No long-lived release token or npm token is required.
4. Protect the default branch and require the `check` job. Configure the
   `github-pages` environment to permit deployment from the default branch.

GitLab:

1. Mirror the GitHub default branch and tags into the GitLab project, preserving
   commit IDs. Configure a GitLab pull mirror, or an existing external mirror
   service. Mirroring credentials are deliberately separate from release jobs;
   these workflows do not push between providers. Ensure mirror updates trigger
   pipelines, and include tags in the mirror configuration.
2. Keep Package Registry and Pages enabled. Allow the built-in job token to
   publish project releases/packages, and protect `v*` tags against independent
   creation. Remove the old tag-creation `GITLAB_TOKEN` variable if unused.
3. The GoReleaser GitLab target is `off-peak.engineer/utilities/patchboard`;
   update `.goreleaser.yml` if moving that backup project.
4. The `pages.publish` syntax requires GitLab 17.9 or newer.

Release and Pages publication happen in the same GitHub workflow as the tested
push; no secondary tag workflow is needed. GitHub tags created with the built-in
token do not start other Actions workflows. See the
[semantic-release GitHub plugin](https://github.com/semantic-release/github),
[GitHub Pages workflow guide](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages),
and [GitLab Pages configuration](https://docs.gitlab.com/ci/yaml/#pagespublish).

## Static board publication

Both pipelines run `npm run build:pages`. Chromium opens `tasks/kanban.html`,
loads the checkout's `tasks/` through a read-only filesystem adapter, and runs
the page's existing parser and renderer. The exporter then saves
`public/index.html`: one fully rendered HTML file with inline styles and task
contents, no JavaScript, no filesystem prompt, and no network dependencies.
Card links point to task text embedded in that file. The editable local board is
unchanged. All configured lanes, including done and anti-feature tasks, are
published. Pages access controls determine who can read this snapshot.

PRs, MRs, and feature branches build preview artifacts but never deploy Pages or
publish releases. Default-branch pushes publish to the provider running the
pipeline. GitHub manual default-branch runs can redeploy Pages without releasing;
GitLab tag pipelines publish backup releases without rolling Pages backward.

## Artifacts and local checks

GoReleaser builds Linux, macOS, and Windows for amd64 and arm64. Archives are
`patchboard_X.Y.Z_OS_ARCH.tar.gz` (Windows: `.zip`) and include the binary,
README, and licenses. `checksums.txt` contains SHA-256 hashes. Snapshot artifacts
expire after one week; published release assets are persistent.

CI uses Go 1.25, the latest stable Node release, npm's committed lockfile, and
checksum-verified GoReleaser 2.18.2. GitHub resolves `current` with
`check-latest: true`; GitLab uses `node:current-bookworm`. The GitHub Actions
themselves use their supported Node 24 runtime, independently of the Node
version used for project commands. Both pipelines print the selected Node and
npm versions. npm enforces the package's minimum Node version (24.10.0).

All direct npm dependencies were checked against the registry's latest releases.
They remain pinned in the lockfile; upgrading dependencies is a reviewed change,
while the Node runtime follows stable releases automatically. For local work,
`nvm install` and `nvm use` read `.nvmrc` and select current Node. Local checks:

```sh
npm ci
npx playwright install chromium
npm test
npm run build:pages
go test -race ./...
go vet ./...
go run ./cmd/patchboard lint
# Linux or macOS archive check:
sh scripts/install-goreleaser.sh
sh scripts/build-archives.sh
```

For a different task directory or output path:
`node scripts/export-board.mjs path/to/tasks path/to/index.html`.
The task directory must contain its own current `kanban.html`.

## Recovery

Before tag creation, fix the failure and rerun the GitHub workflow. If a failure
occurs after semantic-release pushes a tag but before uploading all assets,
inspect that exact tag and GitHub release: semantic-release will not publish the
same version again automatically. Rebuild that tag with
`PATCHBOARD_VERSION=X.Y.Z sh scripts/build-archives.sh`, then recover the missing
release/assets at that tag with the generated notes. Never move or delete a
published tag to force a retry. Fix released code with a new Conventional Commit.
Retry a failed GitLab tag job to recover the backup. Keep release packages out
of registry cleanup rules.

The repository's board config excludes CI caches and `public/` from annotation
scanning so downloaded dependencies and exported task text are not linted as
project source. `npm audit` currently reports three advisories in npm bundled by
semantic-release's transitive npm plugin. That plugin is not enabled here (this
project only publishes Go archives); compatible dependency updates do not yet
resolve those bundled advisories.
