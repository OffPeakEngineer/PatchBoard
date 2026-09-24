# Development

PatchBoard is intentionally straightforward Go. The implementation favors a
small, inspectable filesystem contract over language-specific parsers or a
separate persistence layer.

## Repository Layout

- `cmd/patchboard/main.go` is the process entry point.
- `main.go` preserves the root `go run .` entry point.
- `internal/cli/parse.go` declares and parses the command tree.
- `internal/cli/commands.go` handles command execution, output, and exit
  behavior.
- `internal/tasks` contains configuration, task scanning, annotation scanning,
  lint rules, board operations, repair, and undo behavior.
- `templates` contains embedded task-board assets.
- `tasks` is PatchBoard's dogfood board and durable project-planning record.

## Run Tests

```bash
go test ./...
go test -race ./...
go vet ./...
```

The minimum Go version is 1.25, as declared in `go.mod` and used by GitLab CI.
See [Releases](releases.md) for local archive builds and release checks.

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

Read [CONTRIBUTING.md](../CONTRIBUTING.md) before making contributions. It
defines the project's licensing, authorship, and source-header requirements.
