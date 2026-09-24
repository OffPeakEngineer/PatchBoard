---
id: task-20260714-git-bounded-repository-discovery
owner: andy
tags:
  - patchboard
  - cli
  - git
  - discovery
created: 2026-07-14
---

# Implement explicit and Git-bounded repository discovery

## Problem

Patchboard can select its own dogfood board when invoked from a nested checkout
or submodule, even when the user intended to operate on the containing project.
The executable's location should not influence which repository it manages.

## Proposed behavior

Add a global explicit repository option and make automatic discovery respect the
current Git worktree boundary. Keep repository selection in one shared resolver
used by read-only commands, mutating commands, and `init`.

## Done when

- [x] The reviewed repository resolution precedence is implemented in one place
- [x] `-C`/`--repo` can explicitly select a repository for every subcommand
- [x] Global repository flags work in the documented positions supported by the new
  parser
- [x] Automatic discovery does not cross the current Git worktree boundary
- [x] `init` uses its reviewed default target and does not accidentally initialize a
  nested working directory
- [x] `doctor` reports the resolved repository and task root clearly
- [x] Tests cover a normal root, a nested directory, a nested Git repository, a Git
  submodule, an explicit external path, and a non-Git directory
- [x] Existing commands continue to pass their behavioral tests

## Implementation

- Discovery finds the nearest board but stops at the current `.git` directory or
  gitfile boundary.
- Non-Git boards continue to support upward discovery.
- New non-Git initialization requires `--force`; explicit initialization remains
  available through `-C`, `--repo`, or the legacy positional path.
- Status and doctor identify their resolved repository and task root.

## Depends on

- `task-20260714-repository-resolution-contract`
- `task-20260714-flaggy-cli-migration`
