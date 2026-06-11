---
id: task-20260610-browser-kanban-validation
owner: unassigned
tags:
  - patchboard
  - web
  - compatibility
created: 2026-06-10
---

# Browser kanban enhancements

1) patchboard lint should be able to tell when the template and the 'active template' (at least for kanban.html) drifts, and the 'doctor' or 'fix' will install the new one by copying out that template

2) Currently it allows us to "open a folder", how feasible might it be to look at the URL of the html and deduce that we're in the local filesystem and the taskboard we want to look at is at that directory to at least auto-open the right place from the get-go?

3) Double clicking on a ticket and/or having a document icon that can be clicked that would call the proper "OS exec" to open the markdown file. Leave the file association to the OS/user

4) Offer an "install app" prompt/icon that will keep this project kanban as a little desktop icon app

## Resolution

- `patchboard lint` reports `KANBAN001` when `tasks/kanban.html` is missing or
  differs from the embedded template.
- `patchboard doctor` explains the repair path for kanban drift.
- `patchboard fix` installs or updates `tasks/kanban.html` from
  `templates/kanban.html`.
- The static board shows a local-file hint when opened with `file://`; browsers
  still require an explicit directory picker before granting filesystem access.
- Cards now expose an `Open` file link and double-click behavior for the
  Markdown file. The browser and OS decide whether that opens in-browser or via
  the user's file association.
- The board exposes an install button that calls the browser install prompt when
  available and explains when the browser does not expose that capability.

## Remaining browser reality

Static HTML cannot silently infer or open the surrounding local directory with
write permission. The user must grant access through browser-supported local
filesystem APIs. PWA install is also browser-controlled and may not be available
from a local file URL.
