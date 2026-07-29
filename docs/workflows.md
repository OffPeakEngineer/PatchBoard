# Workflows

PatchBoard provides small commands that fit into local development and CI
without changing its filesystem-first model.

## Diagnose Before Repairing

Use `doctor` for an explanation of board setup and lint findings:

```bash
patchboard doctor
```

Use a dry run before applying safe repairs:

```bash
patchboard fix --dry-run
patchboard fix
```

Repairs are limited to mechanical changes such as missing scaffolding, template
installation, and redundant frontmatter status removal. Decisions such as task
state or filename remain with collaborators.

## Undo with Git

Preview task-tree changes:

```bash
patchboard undo
```

Restore tracked changes and recognized file moves explicitly:

```bash
patchboard undo --apply
```

The command refuses to apply while unrelated untracked task files are present.
This protects notes that Git cannot restore.

## Continuous Integration

`patchboard lint` is the main CI entry point. It returns exit code `1` for lint
errors and supports JSON output for downstream tooling:

```bash
patchboard lint --json
```

Filename conventions can remain warnings or be elevated to errors through board
configuration.

This repository runs GitHub Actions and a parallel GitLab CI pipeline. GitHub-
specific release helpers remain on GitHub; the GitLab pipeline mirrors build,
test, and release behavior where possible. Configure equivalent release secrets
as GitLab CI variables when using those jobs.

## Follow-Up Generation

Future follow-up generation should be explicit: a command may turn selected
lint findings into planning task files when requested. `patchboard lint` should
remain read-only and must not silently mutate a repository.

## Release Support

Release-oriented automation is planned, including workflows for handling done
items during release preparation. Until then, Git history and the board remain
the authoritative inputs rather than generated release state.

GoReleaser builds standalone `patchboard` binaries for Linux, macOS, and
Windows. Release archives contain the applicable license files alongside the
binary. The binary—not a source submodule—is the distributed tool; consuming
repositories retain only their own `tasks/` tree.
