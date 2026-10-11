# Configuration

PatchBoard works without configuration. A repository can opt into local
conventions with `tasks/board.yml`.

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

## Built-in Defaults

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
ignore_dirs:
  - .git
  - node_modules
  - vendor
  - dist
  - build
```

The default annotation markers are listed in [Code Annotations](annotations.md).

When a field is absent or empty, PatchBoard fills it from the built-in defaults.
An explicitly empty `done_states: []` keeps every state nonterminal.

## Lane Configuration Compatibility

Boards can also define columns with `lanes`, using each lane's `path` as its
state folder and `terminal: true` to mark completed or rejected work:

```yaml
lanes:
  - id: anti-feature
    path: -1_anti-feature
    terminal: true
  - id: todo
    path: 1_todo
  - id: doing
    path: 2_doing
  - id: done
    path: 5_done
    terminal: true
```

Lane order controls display order. Explicit `states` take precedence over
`lanes`; explicit `done_states` override the terminal flags. With lane
configuration, no terminal flags means no completed states. The `done` command
recognizes numbered done folders such as `5_done` before falling back to other
terminal states.

`tasks/templates/` holds reusable starters and is excluded from task scanning
unless `templates` is explicitly configured as a state. Lane configuration does
not enable filename lint; `filename.pattern` is only enforced when
`filename.enabled: true` is set, and must then be a regular expression.

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
