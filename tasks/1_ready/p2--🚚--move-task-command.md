---
id: task-20260517-move-task-command
title: Add task movement commands
owner: unassigned
tags:
  - patchboard
  - cli
created: 2026-05-17
---

## Problem

Moving files is the right durable primitive, but the CLI should remove ambiguity
around destination states and any mechanical metadata updates.

## Done when

- `patchboard move <task> <state>` moves a task between configured states
- `patchboard start <task>` moves a task to the configured active lane
- `patchboard done <task>` moves a task to the primary done lane
- Movement preserves the Markdown file as the source of truth
