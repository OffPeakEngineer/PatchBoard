---
id: task-20260518-product-story-docs
owner: andy
tags:
  - docs
  - product
  - story
created: 2026-05-18
---

## Problem

Patchboard's current documentation focuses on CLI commands and developer
workflows. That makes it harder for non-developer stakeholders to understand the
product value, team planning benefits, and why this tool matters.

## Context

The product story should follow the same board model as the dogfood and naming
tasks: `tasks/` is the durable work board, Git history is the audit trail, and
frontmatter is reserved for durable metadata rather than path-derived fields.

## Desired outcome

- A clearer product narrative for non-technical users
- Better documentation of how Patchboard supports team planning and repo-native workflows
- Explicit explanation of `tasks/` as the durable work board and Git history as the audit trail
- Missing docs and feature-story gaps captured as planning items

## Next steps

- Audit README and `tasks/README.md` for user-facing gaps
- Add a "Why Patchboard" section with example use cases
- Add a "Who should use this" section for product, design, and engineering teams
- Surface big misses as additional planning items in `tasks/0_planning/`
