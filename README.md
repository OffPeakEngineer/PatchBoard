# PatchBoard

PatchBoard is a small, opinionated, repo-native task board and annotation
linter.

Tasks are Markdown files. Folders are workflow states. Git history is the audit
trail. The tool makes that system easier to use, but the files remain readable
and meaningful if the binary disappears.

## Why PatchBoard

Product plans and implementation reality often drift apart. A feature may be
marked complete while its code still contains a `TODO`, `FIXME`, or other known
follow-up that never reaches the people planning the work.

PatchBoard keeps both sides close to the repository:

- Product and project collaborators can open `tasks/kanban.html` and use a
  familiar board.
- Engineering collaborators can use Markdown, the CLI, Git, code review, and
  their existing editor.
- Linked code annotations are checked against the same task files that describe
  feature and maintenance work.

PatchBoard is intentionally not a database. The durable record is the file
tree, and the tool helps inspect, repair, and automate it.

The `patchboard` executable is replaceable infrastructure. Its installation
location does not select a project; PatchBoard discovers the board from the
working directory or an explicit `-C`/`--repo` path.

## Who It Is For

PatchBoard fits small product, design, infrastructure, and engineering teams
that want planning context close to implementation context. It is especially
useful when work must remain legible in pull requests, local checkouts, offline
clones, or long-lived repositories.

It is not intended to replace portfolio planning, customer-support queues, or
company-wide reporting. It works best as the repo-local layer of truth: the
concrete plan for what this repository is doing next.

## Install

From a source checkout, install the standalone binary into your configured Go
binary directory:

```bash
go install ./cmd/patchboard
```

Make sure `$(go env GOBIN)`—or `$(go env GOPATH)/bin` when `GOBIN` is empty—is
on your `PATH`. Versioned binaries for Linux, macOS, and Windows (amd64 and arm64)
are published on the [GitLab releases page](https://gitlab.com/off-peak.engineer/utilities/patchboard/-/releases).
The project is currently private, so downloads require project access. See
[release instructions](docs/releases.md) for archive formats and checksums.
Run `patchboard --version` to identify an installed release.

## Quick Start

Initialize a board and inspect it:

```bash
patchboard init
patchboard status
patchboard doctor
```

Create a task, move it into active work, and finish it:

```bash
patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling"
patchboard start fix-login-timeout
patchboard done fix-login-timeout
```

Connect implementation follow-up to a task:

```go
// TODO[task-20260512-auth-timeout]: handle expired refresh token
```

Then inspect or validate annotations:

```bash
patchboard todos
patchboard lint
```

You can also open `tasks/kanban.html` for a visual board. Browsers with local
directory write support can move cards after you choose the project's `tasks/`
directory.

See [Getting Started](docs/getting-started.md) for installation, repository
selection, troubleshooting, development commands, and a fuller first-board
walkthrough.

## Documentation

- [Manifest](docs/manifest.md) — the vision and principles behind PatchBoard
- [Core concepts](docs/concepts.md) — the filesystem contract and design
  principles
- [Getting started](docs/getting-started.md) — initialize and use a board
- [Command reference](docs/commands.md) — CLI commands, JSON output, and exit
  codes
- [Task files](docs/tasks.md) — Markdown format, identity, state, and metadata
- [Code annotations](docs/annotations.md) — marker syntax and lint behavior
- [Configuration](docs/configuration.md) — board config, defaults, and naming
- [Repository resolution](docs/repository-resolution.md) — how PatchBoard
  selects a project and task root
- [Browser kanban](docs/kanban.md) — visual board behavior and limitations
- [Workflows](docs/workflows.md) — doctor, fix, undo, CI, and future releases
- [Development](docs/development.md) — repository architecture and testing
- [Releases](docs/releases.md) — versions, CI, platform downloads, and publication

## Status

PatchBoard is early, but its main filesystem workflow is usable. It can
initialize, inspect, create, move, repair, and validate boards; expose structured
JSON for automation; and provide a browser-based kanban view.

Planned work includes explicit lint follow-up task generation, broader writable
browser support, and done-task cleanup during releases. The dogfood board under
[`tasks/`](tasks/) is the current record of project work.

## Contributing

Read [CONTRIBUTING.md](CONTRIBUTING.md) before submitting changes.

## License

This project is dual-licensed under MIT OR Zlib. See [LICENSE](LICENSE),
[LICENSE.MIT](LICENSE.MIT), and [LICENSE.zlib](LICENSE.zlib).

Copyright (c) 2026 Andrew David LeTourneau
