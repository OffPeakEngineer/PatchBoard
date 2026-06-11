---
id: task-20260610-lint-follow-up-generation
owner: unassigned
tags:
  - patchboard
  - lint
  - automation
created: 2026-06-10
---

# Generate follow-up tasks from lint findings

## Problem

Patchboard can explain lint findings, but teams may want an explicit way to turn
non-blocking findings into planning tasks without making `patchboard lint`
mutate the repo.

## Done when

- A command previews task files that would be generated from selected findings
- Generated tasks use the configured board conventions and task template
- The command requires an explicit apply step before writing files
- Existing lint exit behavior remains strict and predictable

## Open Questions

- Is this describing just a `lint --dryrun`?