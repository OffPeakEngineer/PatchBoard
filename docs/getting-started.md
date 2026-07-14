# Getting Started

PatchBoard works with ordinary files in an existing repository. The examples
below assume the `patchboard` binary is available on your path; during
development, replace `patchboard` with `go run ./cmd/patchboard`.

## Initialize a Board

From the repository root, run:

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

## Use PatchBoard as a Submodule

If PatchBoard is checked out as a submodule, provide the parent repository as
the final argument:

```bash
go -C patchboard run ./cmd/patchboard status ..
go -C patchboard run ./cmd/patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling" ..
go -C patchboard run ./cmd/patchboard move fix-login-timeout doing ..
go -C patchboard run ./cmd/patchboard lint ..
```

See the [command reference](commands.md) for the complete CLI and
[configuration](configuration.md) for repository-specific states and rules.
