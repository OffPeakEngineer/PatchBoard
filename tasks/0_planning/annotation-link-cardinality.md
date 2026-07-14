---
id: task-20260714-annotation-link-cardinality
owner: unassigned
tags:
  - patchboard
  - annotations
  - lint
  - design
created: 2026-07-14
---

# Decide annotation-to-task link cardinality

## Problem

`TODO003` currently enforces one code annotation for each linked task ID. This
provides strict traceability, but a single coherent task may legitimately span
several implementation locations. Requiring a separate task for every location
could create administrative splitting that makes the board less readable.

The opposite extreme also has risk: allowing unlimited links to one task can
turn a broad task into a bucket of unrelated follow-up and make it difficult to
know when the work is actually complete.

## Questions to work through

- Is one annotation per task an essential invariant or an early simplification?
- Should multiple annotations be allowed when they share a task but have unique
  descriptions or locations?
- Would an annotation-level identifier preserve traceability without forcing
  task proliferation?
- Should cardinality be configurable, warned about, or left to board policy?
- How should `patchboard todos`, JSON output, and completion checks represent
  one task with multiple implementation locations?

## Done when

- Real dogfood examples are collected for both one-to-one and one-to-many links
- The intended relationship between tasks and annotations is documented
- `TODO003` is retained, revised, made configurable, or replaced deliberately
- Lint messages explain the chosen rule in terms collaborators can act on
- Tests cover the chosen cardinality and stale-completion behavior
