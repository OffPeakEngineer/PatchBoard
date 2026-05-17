---
id: task-20260517-json-output
title: Add JSON output for automation
owner: unassigned
tags:
  - patchboard
  - automation
created: 2026-05-17
---

## Problem

The human CLI output is useful, but editor integrations, scripts, and CI
annotations need stable structured output.

## Done when

- Status, list, lint, and todos can emit JSON
- JSON includes derived fields without pretending they were hand-authored
- Human output remains the default
