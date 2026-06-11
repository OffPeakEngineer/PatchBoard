---
id: task-20260610-undo-command
owner: unassigned
tags:
  - patchboard
  - cli
  - git
  - repair
created: 2026-06-10
---

# Add an explicit undo command for safe repairs and movement

## Problem

Patchboard uses normal filesystem changes so Git can already undo task moves
and metadata repairs. Some users still need a clear, tool-level command that
explains what can be restored and delegates safely to Git rather than inventing
a second history system.

## Done when

- The command previews the task-board changes it would restore
- The command refuses to overwrite unrelated user edits without confirmation
- The command uses the native OS git commands (not importing the whole golang native git, as a user of a git repo will need to have git installed. I don't need to enforce that the 'patchboard' executable be installed or even built; but git is a hard requirement, since history is integral and part of git that we're using)
- The implementation uses Git history or the worktree index as the source of truth
- Documentation explains when to use Git directly instead
- The kanban.html also supports executing the "undo" commands that are used

## Resolution

- `patchboard undo` previews task-board changes using native `git status`.
- `patchboard undo --apply` restores tracked task-board changes using native
  `git restore`.
- The command refuses unrelated untracked task files instead of deleting them.
- File moves created by Patchboard are handled as the safe pair of a tracked
  deletion plus an untracked task file with the same filename.
- `tasks/kanban.html` supports undoing the most recent browser drag/drop move in
  the current session.
