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

# Define and implement the repository resolution contract

## Problem

Patchboard currently walks upward looking for configuration or a `tasks/`
directory. When Patchboard is dogfooded from a checkout or submodule, more than
one board can be a plausible match and the nearest board may not be the project
the user intended to operate on.

The installed binary should behave consistently regardless of where the binary
itself lives. Automatic discovery and explicit repository selection need one
documented and tested behavior.

## Agreed contract

The intended behavior is documented in `docs/repository-resolution.md`.

- Explicit `-C`/`--repo` selection has highest precedence.
- Legacy positional repository paths remain compatible during migration.
- Supplying both explicit forms is an error.
- Git worktree roots bound automatic discovery.
- Nested repositories and submodules are independent discovery scopes.
- Outside Git, an existing task board remains valid and can be discovered.
- Outside Git with no board, ordinary commands fail clearly and `init` requires
  an explicit force acknowledgment to use the current directory.
- Normal status and doctor output identify the selected roots; short or JSON
  output supports downstream composition.

## Done when

- [x] Resolution precedence is agreed and documented
- [x] Behavior is specified for repository roots, nested directories,
  submodules, worktrees, directories outside Git, and explicitly selected paths
- [x] The default target behavior of `init` is explicitly decided
- [x] The relationship between the legacy positional repo argument and the
  proposed global flag is decided
- [x] A compatibility table captures existing supported invocations that must
  keep working during migration
- [x] Repository resolution is implemented in one shared path
- [x] `-C` and `--repo` work for every command
- [x] Conflicting explicit repository arguments fail clearly
- [x] Automatic discovery respects Git worktree boundaries
- [x] Non-Git discovery and forced initialization match the agreed contract
- [x] Status and doctor report the selected roots as specified
- [x] Tests cover root, nested directory, worktree, nested repository,
  submodule, explicit path, conflicting arguments, and non-Git behavior
- [x] Existing CLI compatibility tests continue to pass
