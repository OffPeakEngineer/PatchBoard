---
id: task-20260513-naming-convention-overhaul
owner: andy
tags:
  - patchboard
  - tasks
  - process
created: 2026-05-13
---

## Problem

The first task naming pass made the board easier to scan, but it also duplicated
too much information across folder names, filenames, root-level configuration,
and frontmatter. That duplication creates maintenance drag for a file-first task
board and makes it harder for Patchboard to distinguish durable task metadata
from fields derived from the file path.

## Decision

Capture the board convention in a single, easy-to-edit task-board config file:
`tasks/board.yml`. Patchboard should stop relying on a root-level
`.patchboard.yaml` for this repo's task-board defaults. Keeping the config under
`tasks/` keeps Patchboard's own metadata close to the board it describes and
avoids adding hidden tool-specific files to the repository root.

The default lane order should be:

- `-1_anti-feature`
- `0_planning`
- `1_ready`
- `2_doing`
- `3_done`

Patchboard may later provide an optional CI or release command that archives or
clears completed tasks so the done lane does not grow without bound.

The current priority-first filename convention is:

```text
p<priority>--<emoji>--<slug>.md
```

Release-candidate detail may be inserted after priority when needed:

```text
p1--rc-1.0.0--🚀--example.md
```

However, the priority prefix and emoji segment have proven awkward for tab
completion and direct shell use. This task should evaluate and document a more
conventional filename format that still preserves the same task information. Any
information removed from the path must either be derivable from configuration or
stored in frontmatter.

Frontmatter should keep durable task metadata such as:

- `id`
- `owner`
- `tags`
- `created`

Frontmatter should not duplicate values that Patchboard can reliably derive from
the task path, such as lane/status, priority, and title. If a field cannot be
derived without losing useful context, it should remain explicit metadata instead
of being silently discarded.

## Done when

- Existing tasks follow the finalized lane and filename convention.
- `tasks/README.md` documents the convention and no longer references stale
  lanes or root-level Patchboard config for this board.
- `tasks/board.yml` captures the desired future Patchboard configuration.
- The frontmatter policy clearly states which fields are durable metadata and
  which fields are derived from the file path.
- Patchboard implementation gaps are captured as follow-up work.
