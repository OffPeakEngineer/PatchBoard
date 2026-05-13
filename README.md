# Patchboard

Patchboard is a small, repo-native task board and annotation linter.

Tasks are Markdown files. Folders are workflow states. Git history is the audit
trail. The tool should make the system easier to use, but the files must stay
readable and meaningful if the binary disappears.

## Status

Patchboard is early. The current useful pieces are:

- `patchboard init`: create the default task board folders
- `patchboard create`: create a task file from command-line fields
- `patchboard lint`: validate task files and linked code annotations
- `patchboard todos`: list code annotations found across the repo

Planned but not built yet:

- JSON output
- task movement commands
- lint follow-up task generation
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
patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling"
patchboard lint
patchboard todos
```

During development, you can run the tool without installing it:

```bash
go run ./cmd/patchboard init
go run ./cmd/patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling"
go run ./cmd/patchboard lint
```

If Patchboard is checked out as a submodule, target the parent repo explicitly:

```bash
go -C patchboard run ./cmd/patchboard init ..
go -C patchboard run ./cmd/patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling" ..
go -C patchboard run ./cmd/patchboard lint ..
go -C patchboard run ./cmd/patchboard todos ..
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
- `TASK005`: filename does not match the configured filename pattern
- `TODO001`: code annotation references a missing task
- `TODO002`: code annotation references a done or archived task
- `TODO003`: duplicate annotation task ID in code

## Configuration

Patchboard works without configuration. Repos can opt into local conventions
with `.patchboard.yaml`, `.patchboard.yml`, or `.patchboard.json` at the repo
root. YAML is nice for project-owned repos because it is comment-friendly:

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
filename:
  enabled: true
  pattern: "^[^-]+--p[0-9]+--[a-z0-9]+(?:-[a-z0-9]+)*\\.md$"
  description: "<type>--pN--slug.md"
  severity: warning
```

JSON is also supported for projects that prefer it:

```json
{
  "task_root": "tasks",
  "states": ["backlog", "ready", "doing", "blocked", "done", "archived"],
  "done_states": ["done", "archived"],
  "filename": {
    "enabled": true,
    "pattern": "^[^-]+--p[0-9]+--[a-z0-9]+(?:-[a-z0-9]+)*\\.md$",
    "description": "<type>--pN--slug.md",
    "severity": "warning"
  }
}
```

Use filename lint as a kindness, not a trap. A project manager who only lives
inside `tasks/` should get a clear message such as
"filename should match <type>--pN--slug.md", not a Go-shaped stack of nonsense.
Set `"severity":
"error"` only when the team wants CI or hooks to enforce the naming convention.

Creation dates, task identity, and richer metadata should live in frontmatter.
Filename conventions are for quick scanning, not as a replacement for the task
record itself.

Future follow-up generation should be explicit, for example a command that turns
lint findings into backlog task files on request. `patchboard lint` should not
silently mutate the repo.

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
