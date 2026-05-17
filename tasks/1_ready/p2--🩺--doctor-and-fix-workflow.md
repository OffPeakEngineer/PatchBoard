---
id: task-20260517-doctor-and-fix-workflow
title: Add doctor and fix workflows
owner: unassigned
tags:
  - patchboard
  - lint
  - repair
created: 2026-05-17
---

## Problem

Patchboard should stay friendly to people who edit files directly. The CLI
needs a safe way to explain and repair mechanical drift without making lint
feel hostile.

## Done when

- `patchboard doctor` reports board health in action-oriented language
- `patchboard fix --dry-run` previews safe mechanical repairs
- `patchboard fix` only performs low-risk file updates
- CI-oriented lint remains strict and predictable
