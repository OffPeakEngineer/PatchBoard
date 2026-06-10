# Patchboard

Patchboard is a small, repo-native task board and annotation linter.

Tasks are Markdown files. Folders are workflow states. Git history is the audit
trail. The tool should make the system easier to use, but the files must stay
readable and meaningful if the binary disappears.

## Status

Patchboard is early. The current useful pieces are:

- `patchboard` or `patchboard status`: summarize board health
- `patchboard init`: create the default task board folders and `tasks/kanban.html`
- `patchboard create`: create a task file from command-line fields
- `patchboard move`: move a task file between configured states
- `patchboard start`: move a task to the active state
- `patchboard done`: move a task to the primary done state
- `patchboard doctor`: explain board setup and lint findings
- `patchboard list`: list tasks by workflow state
- `patchboard lint`: validate task files and linked code annotations
- `patchboard todos`: list code annotations found across the repo
- `--json`: emit structured output for `status`, `list`, `lint`, and `todos`

Planned but not built yet:

- lint follow-up task generation
- broader browser support for writable kanban moves

## Quick Start

Initialize a task board in your repo:

```bash
patchboard init
```

That creates:

```text
tasks/
  README.md
  kanban.html
  backlog/
  ready/
  doing/
  blocked/
  done/
  archived/
```

Then run:

```bash
patchboard
patchboard doctor
patchboard list
patchboard list --json 2_doing
patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling"
patchboard start fix-login-timeout
patchboard done fix-login-timeout
patchboard lint
patchboard lint --json
patchboard todos
```

You can also open `tasks/kanban.html` in a browser for a visual board. Browsers
with local directory write support can move cards between state folders after
you choose the project `tasks/` directory.

During development, you can run the tool without installing it:

```bash
go run ./cmd/patchboard init
go run ./cmd/patchboard status
go run ./cmd/patchboard doctor
go run ./cmd/patchboard list
go run ./cmd/patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling"
go run ./cmd/patchboard move fix-login-timeout doing
go run ./cmd/patchboard lint
```

If Patchboard is checked out as a submodule, target the parent repo explicitly:

```bash
go -C patchboard run ./cmd/patchboard init ..
go -C patchboard run ./cmd/patchboard status ..
go -C patchboard run ./cmd/patchboard doctor ..
go -C patchboard run ./cmd/patchboard list ready ..
go -C patchboard run ./cmd/patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling" ..
go -C patchboard run ./cmd/patchboard move fix-login-timeout doing ..
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
owner: andy
tags:
  - auth
  - bug
created: 2026-05-12
---

# Fix login timeout handling

## Problem

Users can get stuck after their session expires.

## Done when

- Expired sessions redirect cleanly
- Existing session refresh behavior still works
- Regression test added
```

Identity rules:

- Task ID: frontmatter `id`, otherwise the filename slug
- Task title: the first H1 Markdown heading, otherwise frontmatter `title`,
  otherwise the filename slug. New tasks should prefer a Markdown heading.
- Task state: parent folder under `tasks/`

Because the folder is authoritative, frontmatter status is redundant. If a
legacy task includes it, this is invalid:

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
with `tasks/board.yml`, `tasks/board.yaml`, or `tasks/board.json`. Keeping the
config under `tasks/` keeps board metadata near the board it describes:

```yaml
task_root: tasks
states:
  - -1_anti-feature
  - 0_planning
  - 1_ready
  - 2_doing
  - 3_done
done_states:
  - -1_anti-feature
  - 3_done
filename:
  enabled: true
  pattern: "^[a-z0-9]+(?:-[a-z0-9]+)*\\.md$"
  description: "slug.md"
  severity: warning
```

Legacy root config files named `.patchboard.yaml`, `.patchboard.yml`, or
`.patchboard.json` are still supported for existing repos. Board-local config
takes precedence when both are present.

JSON is also supported for projects that prefer it:

```json
{
  "task_root": "tasks",
  "states": ["-1_anti-feature", "0_planning", "1_ready", "2_doing", "3_done"],
  "done_states": ["-1_anti-feature", "3_done"],
  "filename": {
    "enabled": true,
    "pattern": "^[a-z0-9]+(?:-[a-z0-9]+)*\\.md$",
    "description": "slug.md",
    "severity": "warning"
  }
}
```

Use filename lint as a kindness, not a trap. A project manager who only lives
inside `tasks/` should get a clear message such as
"filename should match slug.md", not a Go-shaped stack of nonsense.
Set `"severity":
"error"` only when the team wants CI or hooks to enforce the naming convention.

Creation dates, task identity, ownership, tags, and richer durable metadata
should live in frontmatter. Status comes from the containing folder. Titles can
come from Markdown content or from the filename when the board convention makes
that reliable.

Future follow-up generation should be explicit, for example a command that turns
lint findings into planning task files on request. `patchboard lint` should not
silently mutate the repo.

Exit codes:

- `0`: clean
- `1`: lint errors
- `2`: config, repo, or tooling error

## Built-in Defaults

Patchboard still works without config. The built-in defaults remain generic so
existing repos can opt in without learning a naming scheme first:

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

This project is dual-licensed under MIT OR Zlib. See [LICENSE](LICENSE), [LICENSE.MIT](LICENSE.MIT), and [LICENSE.zlib](LICENSE.zlib).

Copyright (c) 2026 Andrew David LeTourneau
