---
id: task-20260517-centralized-task-database
title: Replace Markdown tasks with a centralized database
owner: unassigned
tags:
  - patchboard
  - anti-feature
created: 2026-05-17
---

## Problem

A database-backed task system would make Patchboard more conventional, but it
would break the repo-native contract that makes the tool useful.

## Done when

- Markdown files remain the durable task representation
- Git history remains the audit trail
- Any future UI or automation treats files as the source of truth
