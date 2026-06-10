---
id: task-20260517-create-defaults-to-configured-state
title: Default create to the first configured state
owner: codex
tags:
  - patchboard
  - cli
  - config
created: 2026-05-17
---

## Problem

`patchboard create` had a hard-coded `backlog` default. That breaks as soon as a
repo dogfoods custom lanes such as `0_backlog`, `1_ready`, and `2_doing`.

## Done when

- The CLI does not force a default state before loading repo config
- Task creation uses the first configured state when no state is supplied
- Existing explicit `--state` behavior still validates against config
