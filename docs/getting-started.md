# Getting Started

PatchBoard works with ordinary files in an existing repository. The executable
does not need to live in that repository and its installation location does not
affect board discovery.

## Install the Binary

From a PatchBoard source checkout, install the standalone executable:

```bash
go install ./cmd/patchboard
```

Go writes it to `GOBIN` when configured, otherwise to `$(go env GOPATH)/bin`.
Add that directory to your shell's `PATH`, then verify the installation:

```bash
patchboard --help
```

Tagged releases are also configured to provide archives for Linux, macOS, and
Windows. Extract the archive for your platform and place the `patchboard` binary
in a directory on `PATH`.

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
  backlog/
  ready/
  doing/
  blocked/
  done/
  archived/
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
patchboard move fix-login-timeout blocked
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
