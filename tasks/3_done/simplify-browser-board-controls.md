---
id: task-20260714-simplify-browser-board-controls
owner: unassigned
tags:
  - patchboard
  - browser
  - kanban
  - usability
created: 2026-07-14
---

# Simplify browser board controls and capability states

## Problem

The browser board presents `Open tasks folder`, `Refresh`, `Undo move`, and
`Install app` together before the page knows which operations are available.
Some controls are irrelevant until a directory is open, while installation is
unlikely to be available when the board is opened as a standalone `file:` URL.
Canceled pickers and unsupported capabilities can also look like failures rather
than normal states.

The interface should stay small and make the next useful action obvious.

## Proposed cleanup

- Make the primary action describe the current state: `Choose tasks folder`,
  `Reconnect tasks folder`, or `Change folder`
- Disable or hide refresh until a board handle is available
- Keep undo hidden until a browser move exists
- Hide the install control unless a real install prompt is available
- Treat picker cancellation as a neutral action rather than an error
- Distinguish read-only, permission-needed, writable, and unsupported states in
  short status text
- Validate that the chosen folder resembles a Patchboard task root and explain
  how to recover when it does not
- Remove small template inconsistencies, including duplicate translation keys
  and product-name casing drift

## Done when

- The empty state has one obvious primary action
- Controls appear or become enabled only when they can do useful work
- Capability and permission states have concise, non-technical messages
- Canceling selection leaves the page stable and usable
- Selecting the wrong directory produces actionable guidance
- English and Spanish messages have matching, non-duplicated keys
- The board remains a single self-contained HTML template with no framework or
  build step

## Depends on

- `task-20260714-browser-capability-and-fallbacks`

## Resolution

Completed 2026-07-29.

- The primary action now changes among choose, reconnect, and change according
  to the current capability and permission state.
- Refresh, undo, forget, change, and install controls appear only when useful.
  Read-only cards are not draggable.
- Canceling the native picker leaves the page unchanged. Wrong directories,
  denied permission, and stale handles produce short recovery guidance.
- The selected directory must contain a usable `board.yml` plus configured
  state folders, or resemble a standard PatchBoard task root.
- English and Spanish translation keys are unique and tested for exact parity.
  PatchBoard casing is consistent in the page.
- The board remains one self-contained HTML template with no framework or build
  step.
