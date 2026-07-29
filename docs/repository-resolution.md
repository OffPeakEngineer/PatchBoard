# Repository Resolution

PatchBoard operates on task files in a project repository. The location of the
PatchBoard executable does not influence which project it selects.

This document defines how commands choose a repository and task root. It is the
behavioral contract for the CLI parser and repository-discovery implementation.

## Terminology

- **Starting directory:** the directory from which automatic discovery begins.
- **Repository root:** the directory PatchBoard treats as the project root.
- **Task root:** the configured task directory beneath the repository root;
  `tasks/` by default.
- **Git worktree root:** the top-level directory of the Git worktree containing
  the starting directory.
- **Explicit repository:** a path supplied with `-C`, `--repo`, or the legacy
  positional repository argument.

## Resolution Precedence

PatchBoard resolves the repository in this order:

1. An explicit `-C` or `--repo` path.
2. The legacy positional repository path, when the command supports it.
3. Automatic discovery from the current working directory.

Supplying both a global repository flag and a positional repository path is an
error. PatchBoard does not guess which explicit value the user intended.

Explicit paths are resolved to absolute paths before board discovery. The
executable's own path is never considered.

## Automatic Discovery Inside Git

When the starting directory belongs to a Git worktree, PatchBoard searches from
that directory upward toward the worktree root.

At each directory it checks, in order:

1. Board-local configuration under the configured task root:
   `tasks/board.yml`, `tasks/board.yaml`, or `tasks/board.json`.
2. Legacy root configuration: `.patchboard.yaml`, `.patchboard.yml`, or
   `.patchboard.json`.
3. An existing `tasks/` directory using built-in defaults.

The nearest matching board wins. Discovery stops at the Git worktree root and
does not inspect its parent directories.

This boundary is important for nested repositories and submodules. Starting
inside a submodule selects a board in that submodule, not a board in its parent
repository. To operate on the parent deliberately, provide `-C` or `--repo`.

Git worktrees follow the same rule: their own reported top-level directory is
the discovery boundary.

## Automatic Discovery Outside Git

A board does not require Git to be readable or usable. Outside a Git worktree,
PatchBoard searches upward for board-local configuration, legacy root
configuration, or an existing `tasks/` directory.

If none is found, ordinary commands return an error explaining that no board or
Git repository could be resolved. They should suggest an explicit repository
path or initialization rather than silently choosing an unrelated directory.

## Initialization

`patchboard init` has slightly different behavior because it creates a board.

- With `-C`, `--repo`, or a legacy positional path, it initializes that explicit
  location.
- Inside Git, it initializes the current Git worktree root by default.
- Outside Git, it initializes an existing discovered board root when one is
  present.
- Outside Git with no existing board, it refuses by default and explains that
  there is no project boundary.
- A force option may explicitly initialize the current directory outside Git.

The force option is an acknowledgment of the target, not permission to
overwrite existing files. Initialization retains its non-overwriting behavior.

## Command Output

Normal `patchboard status` output identifies the selected repository and task
root, similar to how Git reports repository context before describing state.

A short status mode is intended for downstream CLI composition and omits the
contextual prose while retaining predictable task-oriented output. Structured
automation should continue to prefer `--json` where available.

`patchboard doctor` always reports the resolved repository and task root because
repository selection is part of the condition it diagnoses.

## Explicit Repository Flags

`-C` and `--repo` are aliases for the same global option. The option applies to
every command and may appear in any position supported by the CLI parser.

Examples:

```bash
patchboard -C ../project status
patchboard status --repo ../project
patchboard --repo /absolute/project lint
```

The selected directory is the starting point for resolving its PatchBoard
configuration. It is not required to be the final task root itself.

## Compatibility During Migration

Existing invocations remain supported while the global repository option is
introduced:

| Existing form | Compatibility behavior |
| --- | --- |
| `patchboard` | Status for the automatically resolved repository |
| `patchboard status` | Status for the automatically resolved repository |
| `patchboard status PATH` | Preserve as a legacy explicit repository path |
| `patchboard list STATE PATH` | Preserve state followed by repository path |
| `patchboard create [flags] PATH` | Preserve the trailing repository path |
| `patchboard move TASK STATE PATH` | Preserve task, state, and trailing repository path |
| `patchboard start TASK PATH` | Preserve task and trailing repository path |
| `patchboard done TASK PATH` | Preserve task and trailing repository path |
| `patchboard init PATH` | Preserve explicit initialization target |
| `patchboard doctor PATH` | Preserve explicit diagnostic target |
| `patchboard fix [flags] PATH` | Preserve the trailing repository path |
| `patchboard undo [flags] PATH` | Preserve the trailing repository path |
| `patchboard lint [flags] PATH` | Preserve the trailing repository path |
| `patchboard todos [flags] PATH` | Preserve the trailing repository path |

New `-C` and `--repo` forms should be documented as canonical. Legacy
positional paths remain compatible until a separate reviewed deprecation
decision is made.

## Failure Principles

Repository-resolution errors should include:

- the directory from which discovery started;
- whether a Git boundary was found;
- the configuration or task-root names PatchBoard searched for; and
- an actionable next step such as `patchboard init`, `patchboard init --force`,
  or `patchboard -C PATH COMMAND`.

PatchBoard should fail clearly rather than operate on a plausible but unintended
repository.
