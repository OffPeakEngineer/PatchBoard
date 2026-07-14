---
id: task-20260714-standalone-binary-workflow
owner: andy
tags:
  - patchboard
  - distribution
  - documentation
  - workflow
created: 2026-07-14
---

# Make the standalone binary the primary workflow

## Problem

Patchboard was initially described as something a repository could include as a
submodule and run from source. In practice, installing the binary and invoking
it anywhere inside the project is simpler and avoids ambiguity between a
consumer repository's board and Patchboard's own dogfood board.

Only the generated `tasks/` tree needs to be durable in a consuming repository.
The executable should be replaceable infrastructure rather than repository
content.

## Done when

- The quick start presents an installed `patchboard` binary as the normal usage
- Documentation states that executable location does not affect repository
  selection
- Installation guidance covers the supported Go install or release-binary path
- Submodule execution is either removed or moved to a clearly labeled advanced
  or legacy section with an explicit repository target
- Examples consistently demonstrate `-C`/`--repo` when selecting a different
  repository
- Troubleshooting explains how to inspect the repository Patchboard selected
- Release and packaging documentation agree with the standalone-binary model

## Depends on

- `task-20260714-git-bounded-repository-discovery`

## Open questions

- Should source-based submodule execution remain an officially supported
  workflow?
- Is `go install` sufficient as the first documented installation path, or
  should release artifacts be presented first?
