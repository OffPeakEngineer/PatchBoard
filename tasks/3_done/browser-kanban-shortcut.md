---
id: task-20260518-browser-kanban-shortcut
owner: andy
tags:
  - ux
  - docs
  - web
created: 2026-05-18
---

# Build a browser-based kanban shortcut for `tasks/` folder navigation

## Problem

Non-technical PMs need a simple way to view `tasks/` as a traditional board without using the command line or understanding repo internals.

## Desired outcome

- `patchboard init` installs `tasks/kanban.html` from `templates/kanban.html`
- The installed file renders `tasks/` as a kanban-style board without requiring
  Patchboard to be installed on the viewer's system
- Drag-and-drop support moves tasks between states by renaming/moving files when
  the browser grants write access to the selected `tasks/` directory
- A one-click or double-click entry point from the `tasks/` folder or repo root
  that opens the board in a browser
- A durable implementation that keeps Markdown file contents as the source of truth
- Templates are structured for multilingual teams, with visible UI strings
  isolated behind i18n message keys instead of hard-coded throughout behavior

## Notes

The repo now has the first static browser board template. The existing `tasks`
model deliberately keeps UI-specific code at the edges, so this should remain a
project-local lens over Markdown files rather than a second task database.


## Concept v1

Due to browser sandboxing, a plain HTML file cannot silently browse and mutate
the surrounding repository. The board should ask the user to choose the project
`tasks/` directory, then use browser-supported local filesystem access for
read/write operations when available.

Installing the page as an app can make the board easier to launch later, but it
is not the permission model. The permission model is explicit directory access
granted by the user.

## Implementation state

- `templates/kanban.html` exists as a self-contained static board.
- `patchboard init` installs it to `tasks/kanban.html` without overwriting local edits.
- `templates/README.md.tmpl` and `templates/task.md.tmpl` hold generated Markdown
  text that used to live inside Go string literals.
- Browser write behavior still needs hands-on validation across Chrome, Edge,
  Firefox, Safari, macOS, Windows, and Linux.
