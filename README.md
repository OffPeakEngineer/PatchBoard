# Patchboard

Patchboard is a small, repo-native task board and annotation linter.

Tasks are Markdown files. Folders are workflow states. Git history is the audit
trail. The tool should make the system easier to use, but the files must stay
readable and meaningful if the binary disappears.

## Why Patchboard

Patchboard gives a team a task board that lives in the same place as the work.
Instead of copying status between a ticket system, a spreadsheet, and a repo,
the board is ordinary Markdown files under `tasks/`. Moving a task is a file
move. Editing a task is a normal text edit. Reviewing how the plan changed is a
normal Git history question.

That makes the board accessible from two directions:

- Product and project collaborators can open `tasks/kanban.html` and see a
  familiar board.
- Engineering collaborators can use the CLI, shell, Vim, Git, and code review
  without leaving the repo-native workflow.

Patchboard is intentionally not a database. The durable record is the file tree,
and the tool is there to make that tree easier to inspect, repair, and automate.

## Who Should Use This

Patchboard fits small product, design, infrastructure, and engineering teams
that want planning context close to implementation context. It is especially
useful when work needs to remain legible in pull requests, local checkouts,
offline clones, or long-lived repositories.

It is not trying to replace full portfolio planning, customer support queues, or
company-wide reporting systems. It works best as the repo-local layer of truth:
the concrete plan for what this repository is doing next.

## Status

Patchboard is early. The current useful pieces are:

- `patchboard` or `patchboard status`: summarize board health
- `patchboard init`: create the default task board folders and `tasks/kanban.html`
- `patchboard create`: create a task file from command-line fields
- `patchboard move`: move a task file between configured states
- `patchboard start`: move a task to the active state
- `patchboard done`: move a task to the primary done state
- `patchboard doctor`: explain board setup and lint findings
- `patchboard fix`: apply safe mechanical board repairs
- `patchboard undo`: preview or restore task-board changes with Git
- `patchboard list`: list tasks by workflow state
- `patchboard lint`: validate task files and linked code annotations
- `patchboard todos`: list code annotations found across the repo
- `--json`: emit structured output for `status`, `list`, `lint`, and `todos`

CI

This repository runs CI in tandem: existing GitHub Actions workflows remain
under `.github/workflows/`, and a GitLab CI pipeline was added at
`.gitlab-ci.yml` to provide parallel builds on GitLab. Some GitHub-specific
automation (for example, GitHub-only release helpers) are kept on the
GitHub side; the GitLab pipeline mirrors the build/test/release flow where
possible. Set the same secrets in your GitLab project CI variables for
release jobs.

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
  -1_anti-feature/
  0_planning/
  1_ready/
  2_doing/
  3_done/
```

Common commands:

```bash
patchboard
patchboard doctor
patchboard list
patchboard list --json 2_doing
patchboard create --state 1_ready --slug fix-login-timeout --title "Fix login timeout handling"
patchboard start fix-login-timeout
patchboard done fix-login-timeout
patchboard fix --dry-run
patchboard undo
patchboard lint
patchboard lint --json
patchboard todos
```

You can also open `tasks/kanban.html` in a browser for a visual board. Browsers
with local directory write support can move cards between state folders after
you choose the project `tasks/` directory.

Command notes:

- `patchboard` and `patchboard status` summarize board health.
- `doctor` explains setup and lint findings.
- `list [state]` lists tasks, optionally narrowed to one workflow state.
- `create` writes a new task file from command-line fields.
- `move <task> <state>` moves a task between workflow states.
- `start <task>` moves a task to the active state.
- `done <task>` moves a task to the primary done state.
- `fix --dry-run` previews safe mechanical repairs.
- `undo` previews task-board changes that Git can restore.
- `lint` validates task files and linked code annotations.
- `todos` lists code annotations found across the repo.
- `--json` emits structured output for `status`, `list`, `lint`, and `todos`.
- Most commands accept an optional final repo path when Patchboard is run from
  outside the repo that owns `tasks/`.

During development, you can run the tool without installing it. From a repo that
has Patchboard checked out at `./patchboard`, target the parent repo with `..`:

```bash
go -C patchboard run . status ..
go -C patchboard run . doctor ..
go -C patchboard run . list 2_doing ..
go -C patchboard run . create --state 1_ready --slug fix-login-timeout --title "Fix login timeout handling" ..
go -C patchboard run . move fix-login-timeout 2_doing ..
go -C patchboard run . fix --dry-run ..
go -C patchboard run . undo ..
```

If you are working inside the Patchboard source repo itself, the shorter form
works against this repo's own `tasks/` folder:

```bash
go run .
go run . lint
```

You can also build the CLI from a fresh clone with `go build .`.

## Task Files

A task is a Markdown file under one of the state folders:

```text
tasks/2_doing/fix-login-timeout.md
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

Because the folder is authoritative, frontmatter status is redundant. Do not
write this:

```text
tasks/3_done/fix-login-timeout.md
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
- `TASK004`: frontmatter `status` is present
- `TASK005`: filename does not match the configured filename pattern
- `TODO001`: code annotation references a missing task
- `TODO002`: code annotation references a done-state task
- `TODO003`: duplicate annotation task ID in code

## Configuration

Patchboard works without configuration. Repos can opt into local conventions
with `tasks/board.yml`. Keeping the config under `tasks/` keeps board metadata
near the board it describes:

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

`patchboard fix --dry-run` previews safe mechanical repairs. `patchboard fix`
applies only low-risk updates such as creating missing board scaffolding and
removing redundant frontmatter `status` fields. Ambiguous repairs, such as
renaming task files or choosing a state for loose task files, remain doctor
guidance for a human to decide.

`patchboard lint` warns when `tasks/kanban.html` is missing or has drifted from
the bundled template. `patchboard fix` installs or updates that file while still
leaving task Markdown as the source of truth.

`patchboard undo` previews task-board changes that Git can restore.
`patchboard undo --apply` delegates to native `git restore` and refuses to run
while unrelated untracked task files are present. For a browser-only workflow,
`tasks/kanban.html` can undo the most recent drag/drop move during the current
session.

Exit codes:

- `0`: clean
- `1`: lint errors
- `2`: config, repo, or tooling error

## Built-in Defaults

Patchboard still works without config. The built-in defaults match the standard
Patchboard board shape:

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
