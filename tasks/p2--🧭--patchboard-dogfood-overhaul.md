---
id: task-20260513-patchboard-dogfood-overhaul
owner: andy
tags:
  - patchboard
  - dogfood
  - tasks
  - process
created: 2026-05-13
---

## Problem

Patchboard is useful as a repo-native paper trail, but the current dogfood pass
showed places where the tool can fight natural file-first workflows. For now,
the paper trail matters more than perfect lint conformance.

## Lessons

- Moving files between lanes should make status updates easy, forgiving, or
  automatically repairable.
- Lint should distinguish between project-manager-friendly guidance and hard
  failures that block work.
- Templates may need relaxed filename rules or a dedicated template schema.
- The tool should support humans who edit files directly with Vim or an IDE.
- Patchboard may need a command that files its own follow-up tasks when it finds
  non-blocking issues.
- The lane/config model should stay KISS and work from YAML or JSON config.
- Status, priority, and title should derive from path conventions instead of
  being duplicated in frontmatter.
- Repos should be able to define numbered lanes such as `0_backlog`,
  `1_ready`, and terminal lanes such as `x_done` or `-1_closed`.
- Filename lint should support `p<priority>--<emoji>--<slug>.md` and optional
  release-candidate segments such as `p1--rc-1.0.0--🚀--launch.md`.

## Open questions

- Should status derive from folder location, frontmatter, or a reconciled
  combination of both?
- Should `patchboard lint` offer a `--fix` mode for safe frontmatter repairs?
- Which errors should fail CI versus only produce warnings?
- How should moved files preserve history without forcing noisy edits?
- Should Patchboard read `tasks/board.yml`, `.patchboard.yml`, or both?
- How should the tool report derived fields in JSON without pretending they
  were hand-authored metadata?

## Done when

- Patchboard's dogfood requirements are written in its own repo.
- The infrastructure repo can keep using `/tasks` as a paper trail without
  blocking cluster work.
- Lint behavior is friendlier to direct file editing and PM workflows.
