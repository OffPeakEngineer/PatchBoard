# Browser Kanban

`tasks/kanban.html` provides a visual view of a PatchBoard task tree without
requiring the PatchBoard binary or a hosted service.

## Reading a Board

Open the file in a browser and choose the project's `tasks/` directory when
prompted. The page reads configured state folders and renders Markdown task
files as cards.

The browser view is a helper. Task Markdown and folder location remain the
source of truth.

## Moving Cards

Browsers that support local directory write access can move cards between state
folders after the user grants permission. A drag-and-drop move changes the
task's path in the same way as `patchboard move`.

Browser support for local filesystem writes varies. The board remains useful as
a read-only view when writes are unavailable.

## Undo

The page can undo the most recent drag-and-drop move during the current browser
session. It does not replace Git history or the CLI's repository-wide undo
preview.

## Template Updates

PatchBoard bundles the kanban page as a template:

- `patchboard lint` reports `KANBAN001` when the installed file is missing or
  differs from the bundled version.
- `patchboard fix --dry-run` previews an installation or update.
- `patchboard fix` installs the current template.

Updating `kanban.html` is a mechanical repair. It does not alter task Markdown.
