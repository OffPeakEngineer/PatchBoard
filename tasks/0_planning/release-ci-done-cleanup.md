---
id: task-20260610-release-ci-done-cleanup
owner: andy
tags:
  - patchboard
  - release
  - ci
  - safety
created: 2026-06-10
---

# Design release CI cleanup for done tasks

## Problem

Patchboard's done lane is useful as a local paper trail, but it can grow without
bound after work has been released. A CI-oriented command could archive or delete
completed tasks during the release process, but deleting Markdown task records is
destructive and should be designed carefully.

## Context

This repo uses GitLab CI: an MR updates `VERSION` and release notes, a successful
main pipeline creates the tag, and the tag pipeline runs GoReleaser. Archives
and checksums are stored in the Generic Package Registry and linked from the
GitLab release. See `docs/releases.md`. Release jobs currently preserve all tasks.

Current release tooling still commonly builds on Conventional Commits:

- Conventional Commits maps `fix`, `feat`, and breaking changes to SemVer
  release levels and is explicitly designed for changelog/version automation.
- release-please parses Conventional Commit history to create release PRs,
  changelogs, version bumps, and GitHub releases, but is no longer used here.
- semantic-release remains the "take over the whole release workflow" option,
  mostly from the Node ecosystem. This repo currently uses reviewed version
  changes with GoReleaser instead of automatic commit-based version selection.

Patchboard should integrate with that pipeline without owning all release
automation.

## Design questions

- Should the command delete done tasks, move them to an archive lane, or copy
  their summaries into release notes before deletion?
- Should cleanup happen before a release PR, inside the release PR, after a tag,
  or only as a manually approved maintenance command?
- Should cleanup require a clean worktree and explicit `--apply`?
- Should the command refuse to run unless the release tag or release PR exists?
- How should task IDs and task titles map into release notes, if at all?
- Should `-1_anti-feature` be treated as done for lint but excluded from release cleanup?

## Done when

- The release pipeline shape is documented with Patchboard's exact responsibility.
- A proposed command name and mode are chosen, for example
  `patchboard ci cleanup-done --dry-run` and `--apply`.
- The destructive path requires an explicit apply step and refuses dirty or
  ambiguous worktrees.
- The command can be inserted into the existing GitLab CI and GoReleaser
  flow without replacing either tool.
- Documentation explains how Patchboard relates to Conventional Commits,
  release-please, semantic-release, and GoReleaser.
