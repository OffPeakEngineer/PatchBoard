# Getting Started

PatchBoard works with ordinary files in an existing repository. The executable
does not need to live in that repository and its installation location does not
affect board discovery.

## Install the Binary

From a PatchBoard source checkout with Go 1.25 or newer, install the standalone executable:

```bash
go install ./cmd/patchboard
```

Go writes it to `GOBIN` when configured, otherwise to `$(go env GOPATH)/bin`.
Add that directory to your shell's `PATH`, then verify the installation:

```bash
patchboard --help
patchboard --version
```

Official releases are published on
[GitHub](https://github.com/OffPeakEngineer/patchboard/releases). The backup
[GitLab releases page](https://gitlab.com/off-peak.engineer/utilities/patchboard/-/releases)
provides versioned archives for Linux, macOS (`darwin`), and Windows, on amd64
and arm64. Downloads require access to the currently private project. Verify the
archive against `checksums.txt`, extract it, and place `patchboard` (or
`patchboard.exe`) in a directory on `PATH`. See [Releases](releases.md) for details.
Source builds report `dev`; release binaries report the published version.

## Initialize a Board

From anywhere inside the project repository, run:

```bash
patchboard init
```

With the built-in defaults, this creates:

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

Initialization does not overwrite existing board files.

## Inspect Board Health

Run PatchBoard without a command, or use `status` explicitly:

```bash
patchboard
patchboard status
```

For setup guidance and actionable findings, use:

```bash
patchboard doctor
```

## Create and Move a Task

```bash
patchboard create \
  --state ready \
  --slug fix-login-timeout \
  --title "Fix login timeout handling"
```

Move the task into the configured active state and then the primary done state:

```bash
patchboard start fix-login-timeout
patchboard done fix-login-timeout
```

For an explicit destination, use:

```bash
patchboard move fix-login-timeout 0_planning
```

Task IDs, filenames, slugs, and unambiguous task paths can identify a task.

## Link Code Follow-Up

Add a supported annotation to a comment:

```go
// TODO[task-20260512-auth-timeout]: handle expired refresh token
```

List all recognized annotations and validate linked ones:

```bash
patchboard todos
patchboard lint
```

## Use the Browser Board

Open `tasks/kanban.html` in a browser. It can display task files without a
PatchBoard installation. Browsers with local directory write support can also
move cards after you grant access to the project's `tasks/` directory.

## Run from This Source Tree

During development:

```bash
go run ./cmd/patchboard init
go run ./cmd/patchboard status
go run ./cmd/patchboard doctor
go run ./cmd/patchboard list
go run ./cmd/patchboard lint
```

This is a development workflow for PatchBoard itself, not the recommended way
to add PatchBoard to another project. Consuming projects need only the generated
`tasks/` tree in version control.

## Select a Different Repository

Use the global `-C` or `--repo` option when the current directory is not the
project you want to inspect:

```bash
patchboard -C ../another-project status
patchboard lint --repo /absolute/path/to/project
```

Legacy trailing repository arguments remain compatible, but the global option
is the canonical form for new scripts and documentation.

## Troubleshoot Repository Selection

Normal `status` output identifies the repository and task root PatchBoard
selected:

```text
Patchboard status
Repo: /path/to/project
Task root: /path/to/project/tasks
```

`patchboard doctor` reports the same paths alongside setup findings. If the
wrong project is selected, run the command with `-C PATH` rather than moving or
reinstalling the executable.

Automatic discovery stops at the current Git worktree boundary. A nested Git
repository or submodule is therefore its own discovery scope.

## Legacy Submodule Execution

Embedding PatchBoard's source as a submodule is no longer the recommended
installation model. Existing checkouts can still run it explicitly while they
migrate to an installed binary:

```bash
go -C patchboard run ./cmd/patchboard --repo .. status
go -C patchboard run ./cmd/patchboard --repo .. lint
```

The explicit repository option is required because running from the PatchBoard
checkout would otherwise select PatchBoard's own dogfood board.

See the [command reference](commands.md) for the complete CLI and
[configuration](configuration.md) for repository-specific states and rules.
