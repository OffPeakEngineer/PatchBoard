# Command Reference

Commands accept an optional repository path. When omitted, PatchBoard starts
from the current directory and searches upward for a board.

## Global Repository Selection

Every command accepts `-C PATH` or `--repo PATH`. The option can appear before
or after the subcommand:

```bash
patchboard -C ../project status
patchboard lint --repo ../project
```

The global option is preferred over the compatible legacy trailing path.
Supplying both forms is an error.

## Board Views

### `patchboard` and `patchboard status`

Summarize task counts, active work, annotations, and lint health.

```bash
patchboard status
patchboard status --short
patchboard status --json
```

### `patchboard list`

List tasks, optionally limited to a configured state.

```bash
patchboard list
patchboard list ready
patchboard list --json 2_doing
```

### `patchboard todos`

List code annotations recognized across the repository.

```bash
patchboard todos
patchboard todos --json
```

## Board Changes

### `patchboard init`

Create missing default or configured board scaffolding without overwriting
existing files.

```bash
patchboard init
patchboard init --force
```

Outside Git with no existing board, `--force` explicitly initializes the
current directory. Explicit `-C`/`--repo` targets do not require this
acknowledgment.

### `patchboard create`

Create a task from command-line fields. The title is required. State defaults
to the board's preferred planning state, and the slug can be derived from the
title.

```bash
patchboard create --state ready --slug fix-login-timeout --title "Fix login timeout handling"
patchboard create --title "Investigate cache misses" --owner andy --tags performance,backend
```

### `patchboard move`

Move one task to an explicit configured state.

```bash
patchboard move fix-login-timeout doing
```

### `patchboard start`

Move one task to the board's inferred active state.

```bash
patchboard start fix-login-timeout
```

### `patchboard done`

Move one task to the board's primary configured done state.

```bash
patchboard done fix-login-timeout
```

## Validation and Repair

### `patchboard lint`

Validate task files, linked annotations, and the installed browser-board
template.

```bash
patchboard lint
patchboard lint --json
```

### `patchboard doctor`

Explain board setup problems and lint findings with suggested actions.

```bash
patchboard doctor
```

### `patchboard fix`

Apply safe mechanical repairs, including missing scaffolding, the current
kanban template, and removal of redundant frontmatter status fields.

```bash
patchboard fix --dry-run
patchboard fix
```

Ambiguous changes, such as choosing a state for a loose task or renaming a task
file, remain human decisions.

### `patchboard undo`

Preview task-board changes that Git can restore. Apply the restore explicitly:

```bash
patchboard undo
patchboard undo --apply
```

Undo refuses to apply while unrelated untracked task files are present.

## JSON Output

`status`, `list`, `lint`, and `todos` support `--json` for scripts and CI. Human
output is the default.

## Exit Codes

- `0`: the command completed successfully; lint found no errors
- `1`: lint errors were found, or a guarded operation was refused
- `2`: configuration, repository, arguments, or required tooling failed
