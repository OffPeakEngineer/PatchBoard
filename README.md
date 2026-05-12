# Patchboard

Patchboard is a small, repo-native task board and annotation linter.

Tasks are Markdown files. Folders are workflow states. Git history is the audit
trail. The tool should make the system easier to use, but the files must stay
readable and meaningful if the binary disappears.

## Status

Patchboard is early. The current useful pieces are:

- `patchboard init`: create the default task board folders
- `patchboard lint`: validate task files and linked code annotations
- `patchboard todos`: list code annotations found across the repo

Planned but not built yet:

- config file support
- JSON output
- task creation and movement commands
- a local web UI

## Quick Start

Initialize a task board in your repo:

```bash
patchboard init
```

That creates:

```text
tasks/
  README.md
  backlog/
  ready/
  doing/
  blocked/
  done/
  archived/
```

Then run:

```bash
patchboard lint
patchboard todos
```

During development, you can run the tool without installing it:

```bash
go run ./cmd/patchboard init
go run ./cmd/patchboard lint
```

If Patchboard is checked out as a submodule, target the parent repo explicitly:

```bash
go -C patchboard run ./cmd/patchboard init ..
go -C patchboard run ./cmd/patchboard lint
go -C patchboard run ./cmd/patchboard todos
```

## Task Files

A task is a Markdown file under one of the state folders:

```text
tasks/doing/fix-login-timeout.md
```

The parent folder is the authoritative workflow state. Frontmatter is optional,
but useful:

```markdown
---
id: task-20260512-auth-timeout
title: Fix login timeout handling
status: doing
priority: medium
owner: andy
tags:
  - auth
  - bug
created: 2026-05-12
---

## Problem

Users can get stuck after their session expires.

## Done when

- Expired sessions redirect cleanly
- Existing session refresh behavior still works
- Regression test added
```

Identity rules:

- Task ID: frontmatter `id`, otherwise the filename slug
- Task title: frontmatter `title`, otherwise the first Markdown heading,
  otherwise the filename slug
- Task state: parent folder under `tasks/`

Because the folder is authoritative, this is invalid:

```text
tasks/done/fix-login-timeout.md
```

```yaml
status: doing
```

## Code Annotations

Patchboard scans text files for annotation comments. It is not tied to one
language, so these forms are all valid:

```go
// TODO[task-20260512-auth-timeout]: handle expired refresh token
```

```sh
# FIXME[task-20260512-auth-timeout]: shell scripts need the same behavior
```

```html
<!-- WARN[task-20260512-auth-timeout]: this flow is stale -->
```

Use square brackets for task links:

```text
MARKER[task-id]: message
```

Use parentheses for a short owner or tag:

```text
TODO (dave): follow up
WARN (@andy): check before deploy
```

Older `TODO(task-id): message` style is accepted when the value looks like an
ID, but new linked annotations should use brackets. That keeps task IDs distinct
from owners.

Default markers:

- `TODO`
- `FIXME`
- `XXX`
- `WARN`
- `WARNING`
- `BUG`
- `HACK`
- `NOTE`
- `REVIEW`
- `OPTIMIZE`
- `PERF`
- `SECURITY`
- `DEPRECATED`
- `TEMP`
- `TBD`
- `TASK`

Unlinked annotations are listed by `patchboard todos`, but they do not fail
`patchboard lint`. Linked annotations are checked against task IDs.

## Lint Rules

- `TASK001`: unknown state folder
- `TASK002`: duplicate task ID
- `TASK003`: missing title
- `TASK004`: frontmatter `status` does not match the folder
- `TODO001`: code annotation references a missing task
- `TODO002`: code annotation references a done or archived task
- `TODO003`: duplicate annotation task ID in code

Exit codes:

- `0`: clean
- `1`: lint errors
- `2`: config, repo, or tooling error

## Defaults

Patchboard currently works without config:

```yaml
task_root: tasks
states:
  - backlog
  - ready
  - doing
  - blocked
  - done
  - archived
done_states:
  - done
  - archived
ignore_dirs:
  - .git
  - node_modules
  - vendor
  - dist
  - build
```

Files larger than 2 MiB are skipped, and obvious binary files are skipped by
checking for a NUL byte near the start of the file.

## Development

Run tests:

```bash
go test ./...
```

Run the CLI against a repo that contains `tasks/`:

```bash
go run ./cmd/patchboard init
go run ./cmd/patchboard lint
go run ./cmd/patchboard todos
```

The implementation is intentionally boring Go:

- `cmd/patchboard` handles command dispatch and process exit codes.
- `internal/tasks` contains the task model, scanners, and lint rules.
- The scanner uses filesystem walks and regular expressions rather than
  language-specific parsers, because annotations should work across many file
  types.

When adding features, keep the filesystem contract intact: Markdown tasks and
Git history are the durable system. Patchboard is a helper, not the database.

## License

MIT. See [LICENSE](LICENSE).
