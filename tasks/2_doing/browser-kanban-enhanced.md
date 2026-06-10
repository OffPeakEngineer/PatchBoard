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