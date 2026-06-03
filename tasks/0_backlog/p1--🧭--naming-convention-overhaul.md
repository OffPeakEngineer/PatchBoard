

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

The first task naming pass worked, but it duplicated too much information
between folder names, filenames, and frontmatter. That creates maintenance drag
for a file-first task board.

## Decision

- Lanes are prefixed for scan order: `0_backlog`, `1_ready`, `2_doing`,
  `3_blocked`, `x_done`, and `-1_closed`.
- Filenames are priority-first: `p<priority>--<emoji>--<slug>.md`.
- Release-candidate detail may be inserted after priority:
  `p1--rc-1.0.0--🚀--example.md`.
- Frontmatter keeps durable metadata such as `id`, `owner`, `tags`, and
  `created`.
- Frontmatter does not repeat `title`, `status`, or `priority`; those derive
  from the file path.

## Done when

- Existing tasks follow the new lane and filename convention.
- `tasks/README.md` documents the convention.
- `tasks/board.yml` captures the desired future Patchboard config.
- Patchboard implementation gaps are captured as follow-up work.
