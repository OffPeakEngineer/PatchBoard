---
id: task-20260518-browser-kanban-shortcut
title: Build a browser-based kanban shortcut for `tasks/` folder navigation
status: backlog
priority: p3
owner: andy
tags:
  - ux
  - docs
  - web
---

## Problem
Non-technical PMs need a simple way to view `tasks/` as a traditional board without using the command line or understanding repo internals.

## Desired outcome
- A local web app or shortcut that renders `tasks/` as a kanban-style board
- Drag-and-drop support for moving tasks between states by renaming/moving files
- A one-click (or doulbe, if shortcut or app or script file) entry point from the `tasks/` folder or repo root that opens the board in a browser
- A durable implementation that keeps Markdown file contents as the source of truth

## Notes
The repo currently has no built-in web UI, and the existing `tasks` model deliberately keeps UI/UI-specific code at the edges.
A simple web app is the lowest bar and fits the current architecture better than a shell extension for every window manager.
