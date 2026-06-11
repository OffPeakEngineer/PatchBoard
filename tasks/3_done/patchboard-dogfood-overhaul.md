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
showed places where the tool can fight natural file-first workflows. Patchboard
should make the repository task board easier to maintain without becoming a
strict metadata gatekeeper.

## Principles

- The Markdown file remains the durable task record.
- Folder location is the source of truth for lane/status.
- Git history is the audit trail for movement, renames, and task evolution.
- Humans should be able to edit tasks directly in a shell, Vim, or an IDE.
- Lint should distinguish advisory guidance from hard failures.
- Hard failures should be reserved for states Patchboard cannot reconcile or
  explain clearly.
- Project-specific board conventions should live near the board they describe.

## Decisions

The naming convention overhaul is the authoritative decision record for board
layout, filename conventions, config location, and frontmatter policy. For this
repo, Patchboard should use `tasks/board.yml` as the board config and the
default lane set should be:

- `-1_anti-feature`
- `0_planning`
- `1_ready`
- `2_doing`
- `3_done`

Status and title should be derived from path conventions when they can be
derived without losing useful context. Durable metadata belongs in frontmatter.

The broader product may later support additional config filenames or formats,
but this repo should use one canonical board config until that flexibility is
needed and specified.

Patchboard should eventually support repo-defined lane templates and filename
templates using a simple human-readable format. The template language should be
easy to inspect and edit; it should not require users to understand regular
expressions for common board conventions.

## Resolution

The dogfood pass converged on a file-first board that is now reflected in code,
docs, and Patchboard's own `tasks/` tree:

- Board-local config lives in `tasks/board.yml`.
- Task files use shell-friendly `slug.md` names.
- Durable metadata stays in frontmatter.
- Lane/status comes from the containing folder.
- `patchboard move`, `patchboard start`, and `patchboard done` perform
  filesystem movement without mutating task content.
- `patchboard doctor` explains board health, and `patchboard fix` applies only
  low-risk mechanical repairs.
- `patchboard init` installs editable project-local templates, including the
  zero-install `tasks/kanban.html` board lens.

## Follow-up Work

- Browser write behavior is captured in `browser-kanban-validation`.
- Lint follow-up task generation is captured in `lint-follow-up-generation`.
- Git-backed undo behavior is captured in `undo-command`.
- Human-readable convention templates are captured in `board-template-language`.

## Done when

- Patchboard's dogfood requirements are written in its own repo.
- The infrastructure repo can keep using `/tasks` as a paper trail without
  blocking cluster work.
- The naming convention, JSON output, and product-story tasks all point back to
  the same board model.
- Lint behavior is friendlier to direct file editing and planning workflows.
