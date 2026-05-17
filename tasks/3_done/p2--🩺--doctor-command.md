---
id: task-20260517-doctor-command
title: Add a read-only doctor command
owner: codex
tags:
  - patchboard
  - cli
  - dogfood
created: 2026-05-17
---

## Problem

Patchboard needs a friendlier command for explaining board setup problems,
especially while its own repository is being bootstrapped.

## Done when

- `patchboard doctor` prints repo and task-root context
- Missing configured state folders are reported with next actions
- Root-level task files are reported because they are not assigned to a state
- Existing lint findings are shown with action-oriented guidance
