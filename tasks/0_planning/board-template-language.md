---
id: task-20260610-board-template-language
owner: unassigned
tags:
  - patchboard
  - config
  - templates
created: 2026-06-10
---

# Define human-readable board and filename templates

## Problem

Patchboard supports configured states and filename lint patterns, but regular
expressions are not the right interface for most teams. A board convention
should be easy to inspect, explain, and edit without knowing regex syntax.

## Done when

- Common lane sets can be declared without regular expressions
- Common filename conventions can be declared with named fields
- Existing regex-based filename lint remains supported for advanced repos
- `doctor` explains template-derived conventions in human language

## Open Questions

- Not sure I know the full problem here; but I think instead of supporting 