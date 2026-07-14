# Development

PatchBoard is intentionally straightforward Go. The implementation favors a
small, inspectable filesystem contract over language-specific parsers or a
separate persistence layer.

## Repository Layout

- `cmd/patchboard` handles command dispatch, flags, human output, JSON output,
  and process exit codes.
- `internal/tasks` contains configuration, task scanning, annotation scanning,
  lint rules, board operations, repair, and undo behavior.
- `templates` contains embedded task-board assets.
- `tasks` is PatchBoard's dogfood board and durable project-planning record.

## Run Tests

```bash
go test ./...
```

## Run the CLI from Source

```bash
go run ./cmd/patchboard status
go run ./cmd/patchboard doctor
go run ./cmd/patchboard lint
go run ./cmd/patchboard todos
```

## Design Constraint

When adding features, keep the filesystem contract intact: Markdown tasks and
Git history are the durable system. PatchBoard is a helper, not the database.

The annotation scanner uses filesystem walks and regular expressions so it can
work across file types. Changes to scanning should preserve predictable
language-independent behavior and should be covered with fixture-based tests.

## Project Policies

Read [CONTRIBUTING.md](../CONTRIBUTING.md) and [AGENTS.md](../AGENTS.md) before
making contributions. They define the project's licensing, authorship, source
header, and AI-use requirements.
