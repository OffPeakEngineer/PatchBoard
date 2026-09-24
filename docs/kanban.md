# Browser Kanban

`tasks/kanban.html` provides a visual view of a PatchBoard task tree without
requiring the PatchBoard binary or a hosted service.

## Reading a Board

Open the file in a browser and choose the project's `tasks/` directory. The page
checks for `board.yml` and configured state folders before it accepts the
selection, then renders Markdown task files as cards.

The browser view is a helper. Task Markdown and folder location remain the
source of truth.

## Remembering a Folder

After a successful selection, the page stores the browser's directory handle in
IndexedDB. It does not store task contents or maintain a second copy of the
board.

On a later visit, the page calls `queryPermission()` on the remembered handle:

- Granted read/write access loads the board immediately.
- Granted read-only access loads the board without drag-and-drop and offers
  `Reconnect tasks folder`.
- Permission that must be requested again is requested only after the user
  chooses `Reconnect tasks folder`.
- A stale, moved, deleted, or revoked handle keeps `Change folder` and
  `Forget folder` available as recovery actions.

`Change folder` opens the browser picker for another project. `Forget folder`
removes only the saved handle from IndexedDB.

## Moving Cards

Browsers that support local directory write access can move cards between state
folders after the user grants permission. A drag-and-drop move changes the
task's path in the same way as `patchboard move`.

Browser support for local filesystem writes varies. The board remains useful as
a read-only view when writes are unavailable.

## Capability and Control States

| State | Primary action | Other controls | Board behavior |
| --- | --- | --- | --- |
| No folder | Choose tasks folder | Language only | Empty |
| Permission needed | Reconnect tasks folder | Change, forget | Not loaded |
| Read-only | Reconnect tasks folder | Change, refresh, forget | Cards are not draggable |
| Writable | Change folder | Refresh, forget | Drag-and-drop enabled |
| Unsupported | None | Language only | Explanation points to a compatible browser or the CLI |

Refresh is unavailable until a usable handle exists. Browser undo appears only
after a move in the current session. Install appears only when the browser emits
a real app-install prompt. Canceling the folder picker is a neutral action and
does not replace the current board or status.

`showDirectoryPicker()` is not a broadly available web standard. The picker
requires a user gesture, and its stable `patchboard-tasks` ID lets supporting
browsers reuse the picker location independently of the saved directory handle.
See the [API contract and compatibility notes on
MDN](https://developer.mozilla.org/en-US/docs/Web/API/Window/showDirectoryPicker)
and [Chrome's permission guidance](https://developer.chrome.com/docs/capabilities/web-apis/file-system-access).

The release check on 2026-07-29 used Google Chrome 150 on macOS:

| Surface | Result |
| --- | --- |
| Standalone `file:` page | `showDirectoryPicker` exposed; one-action empty state rendered |
| Empty controls | Refresh, undo, forget, change, and install hidden |
| Script validation | Inline JavaScript parsed successfully |
| Picker gesture and OS selection | Manual check required; headless Chrome cannot complete the native directory dialog |
| Firefox, Safari, and non-Chromium browsers | Page detects a missing API and presents the unsupported explanation; confirm behavior on the target browser version |

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

## Repeatable Verification

1. Open `tasks/kanban.html` from the filesystem in a supported Chromium browser.
2. Confirm the empty state shows only `Choose tasks folder` plus language.
3. Cancel the picker and confirm the page remains in the same usable state.
4. Select a non-task folder and confirm the recovery message mentions
   `board.yml` and state folders.
5. Select the project `tasks/` directory, reload, and verify immediate load or
   the single `Reconnect tasks folder` action according to browser permission.
6. Revoke access, reload, and verify reconnect, change, and forget recovery.
7. Move one card, verify undo appears, then undo the move.
8. Run `go test ./...` to check translation key parity and the template's
   remembered-folder contracts.
