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

## Follow-up Work

- Implement or document the finalized naming convention from the naming task.
- Decide how permissive filename lint should be for templates and transitional
  task files.
- Define which lint findings are warnings, which are errors, and which can be
  repaired automatically.
- Decide whether an explicit undo command should wrap Git restoration for safe
  metadata and movement repairs.
- Ensure JSON output reports derived fields without presenting them as
  hand-authored frontmatter.
- Capture non-blocking findings as follow-up tasks when Patchboard can identify
  useful next steps.

## Done when

- Patchboard's dogfood requirements are written in its own repo.
- The infrastructure repo can keep using `/tasks` as a paper trail without
  blocking cluster work.
- The naming convention, JSON output, and product-story tasks all point back to
  the same board model.
- Lint behavior is friendlier to direct file editing and planning workflows.
