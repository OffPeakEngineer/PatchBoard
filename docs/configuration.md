# Configuration

PatchBoard works without configuration. A repository can opt into local
conventions with `tasks/board.yml`, `tasks/board.yaml`, or `tasks/board.json`.

Keeping configuration under `tasks/` places board metadata beside the board it
describes.

## Example

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

JSON configuration is also supported:

```json
{
  "task_root": "tasks",
  "states": ["-1_anti-feature", "0_planning", "1_ready", "2_doing", "3_done"],
  "done_states": ["-1_anti-feature", "3_done"],
  "filename": {
    "enabled": true,
    "pattern": "^[a-z0-9]+(?:-[a-z0-9]+)*\\.md$",
    "description": "slug.md",
    "severity": "warning"
  }
}
```

## Built-in Defaults

```yaml
task_root: tasks
states:
  - backlog
  - ready
  - doing
  - blocked
  - done
  - archived
done_states:
  - done
  - archived
ignore_dirs:
  - .git
  - node_modules
  - vendor
  - dist
  - build
```

The default annotation markers are listed in [Code Annotations](annotations.md).

## Configuration Precedence

Legacy root files named `.patchboard.yaml`, `.patchboard.yml`, or
`.patchboard.json` remain supported. A board-local config takes precedence when
both forms are present.

When a field is absent or empty, PatchBoard fills it from the built-in defaults.

## Filename Rules

Filename lint can be disabled or configured with a regular expression,
human-readable description, and severity.

Use `warning` when the rule is guidance. Use `error` only when the team wants CI
or hooks to enforce it. The description should explain the expected convention
without exposing collaborators to regular-expression details.

## State Semantics

The order of `states` controls display order and helps PatchBoard infer planning
and active states for convenience commands. `done_states` determines which
tasks are considered completed when validating linked annotations. The first
applicable done state is used by `patchboard done`.
