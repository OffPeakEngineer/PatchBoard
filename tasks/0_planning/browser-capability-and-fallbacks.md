---
id: task-20260714-browser-capability-and-fallbacks
owner: unassigned
tags:
  - patchboard
  - browser
  - kanban
  - compatibility
created: 2026-07-14
---

# Define browser board capabilities and fallbacks

## Problem

The browser kanban is an important surface for collaborators who do not use the
CLI, but local filesystem write APIs are not available consistently across
browsers. The board can therefore be writable in one environment and read-only
in another without a sufficiently explicit capability contract.

PatchBoard should keep the browser useful without implying that every supported
browser can safely move files.

## Questions to work through

- Which browsers and filesystem capabilities are officially supported?
- How should the board explain read-only mode and unavailable write access?
- Should unsupported environments load task data through a generated snapshot,
  directory selection, or another non-authoritative adapter?
- Which operations must remain CLI-only because browser filesystem behavior is
  too limited or risky?
- How should browser moves and session undo communicate their relationship to
  Git-backed history and CLI undo?

## Done when

- A small capability matrix documents tested read and write behavior
- The board detects its available capabilities and communicates them clearly
- Read-only use remains useful when directory writes are unavailable
- Writable moves require explicit permission and preserve the filesystem-as-
  truth contract
- Browser undo boundaries are documented and visually distinguishable from Git
  restore behavior
- Compatibility behavior has repeatable manual or automated verification steps
