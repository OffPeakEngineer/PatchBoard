---
id: task-20260517-json-output
owner: unassigned
tags:
  - patchboard
  - automation
created: 2026-05-17
---

## Problem

The human CLI output is useful, but editor integrations, scripts, and CI
annotations need stable structured output.

## Context

This work depends on the board model settled by the naming convention overhaul.
Patchboard needs to expose both durable frontmatter and path-derived fields, but
the JSON shape must make that distinction clear.

## Done when

- Status, list, lint, and todos can emit JSON
- JSON includes derived fields such as lane/status, priority, and title without
  pretending they were hand-authored frontmatter
- JSON includes enough source information for callers to tell whether a field
  came from frontmatter, the file path, or board configuration
- Human output remains the default
