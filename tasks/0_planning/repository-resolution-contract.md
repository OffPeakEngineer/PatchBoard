---
id: task-20260714-repository-resolution-contract
owner: andy
tags:
  - patchboard
  - cli
  - discovery
  - design
created: 2026-07-14
---

# Define the repository resolution contract

## Problem

Patchboard currently walks upward looking for configuration or a `tasks/`
directory. When Patchboard is dogfooded from a checkout or submodule, more than
one board can be a plausible match and the nearest board may not be the project
the user intended to operate on.

The installed binary should behave consistently regardless of where the binary
itself lives. Automatic discovery and explicit repository selection need a
small, documented precedence contract before their implementation changes.

## Proposed contract

1. An explicit repository argument or global `-C`/`--repo` flag wins.
2. Otherwise, Patchboard identifies the Git worktree containing the current
   directory.
3. Discovery searches from the current directory toward that worktree root for
   board-local configuration, legacy root configuration, or `tasks/`.
4. Discovery does not continue above the current Git worktree root.
5. `init` defaults to the Git worktree root when no explicit target is given.
6. Diagnostic output can show which repository and task root were selected.

## Done when

- Resolution precedence is agreed and documented
- Behavior is specified for repository roots, nested directories, submodules,
  worktrees, directories outside Git, and explicitly selected paths
- The default target behavior of `init` is explicitly decided
- The relationship between the legacy positional repo argument and the proposed
  global flag is decided
- A compatibility table captures existing supported invocations that must keep
  working during migration

## Open questions

- Should a directory outside Git still use upward `tasks/` discovery, operate
  only on the current directory, or require an explicit target?
- Should both `-C` and `--repo` be supported as names for the same option?
- Should `status` always display the selected root, or only under `--verbose`?
