---
id: task-20260517-default-status-and-list
title: Add default status and task listing commands
owner: codex
tags:
  - patchboard
  - cli
  - dogfood
created: 2026-05-17
---

## Problem

Patchboard required users to know which command to run before it could help.
That made the board easy to ignore even when the Markdown workflow was useful.

## Done when

- `patchboard` prints a useful board status summary
- `patchboard status` is available as an explicit command
- `patchboard list` prints tasks grouped by configured workflow state
- `patchboard list <state>` filters to one configured state
