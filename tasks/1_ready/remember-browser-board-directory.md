---
id: task-20260714-remember-browser-board-directory
owner: unassigned
tags:
  - patchboard
  - browser
  - kanban
  - usability
created: 2026-07-14
---

# Remember the browser board directory

## Problem

Opening `tasks/kanban.html` from the filesystem does not give the page access to
its parent directory. The user must navigate to the repository to open the HTML
file and then navigate to the same place again in a directory picker before the
board can load.

Browser security requires an explicit user action before granting filesystem
access, so the page cannot silently infer or open its own parent directory. It
can potentially avoid repeating the directory picker after the first successful
selection by storing the granted directory handle in IndexedDB and checking its
permission on later visits.

## Proposed experience

1. The first visit asks the user to choose the `tasks/` directory.
2. The selected `FileSystemDirectoryHandle` is stored locally in IndexedDB.
3. Later visits attempt to restore that handle and query its permission.
4. If permission remains granted, the board loads immediately.
5. If permission must be requested again, the page offers one clear
   `Reconnect tasks folder` action rather than opening an unrelated picker.
6. `Change folder` remains available for switching projects or recovering from
   a stale handle.

## Done when

- The current behavior of the picker `id` is tested for `file:` pages in
  supported browsers before adding storage
- A selected directory handle can be saved and restored without storing task
  contents or introducing another source of truth
- Permission is checked with `queryPermission()` and requested only from a user
  action when necessary
- Returning users avoid directory navigation when the browser permits it
- Stale, moved, revoked, or deleted directory handles fail with a clear recovery
  action
- The page provides an explicit way to forget or change the selected directory
- Unsupported browsers retain a clear read-only or unsupported explanation

## Depends on

- `task-20260714-browser-capability-and-fallbacks`

## Constraint

The initial permission gesture cannot be removed for a standalone local HTML
file. The goal is to make selection a one-time or one-click setup, not to bypass
the browser's filesystem security model.
