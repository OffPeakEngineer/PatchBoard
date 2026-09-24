# Code Annotations

PatchBoard scans text files for configured annotation markers. An annotation can
link implementation follow-up to a task, identify a short owner, or remain an
unlinked inventory item.

## Linked Annotations

Use square brackets for a task ID:

```go
// TODO[task-20260512-auth-timeout]: handle expired refresh token
```

```sh
# FIXME[task-20260512-auth-timeout]: shell scripts need the same behavior
```

```html
<!-- WARN[task-20260512-auth-timeout]: this flow is stale -->
```

The general form is:

```text
MARKER[task-id]: message
```

Older `TODO(task-id): message` annotations are accepted when the value looks
like an ID, but new task links should use brackets.

## Owners and Unlinked Notes

Use parentheses for a short owner or tag:

```text
TODO (dave): follow up
WARN (@andy): check before deploy
```

Markers without a task or owner are also useful:

```text
TODO: investigate this branch
```

Unlinked annotations appear in `patchboard todos`, but they do not fail lint.
Linked annotations claim a relationship and are therefore validated.

## Default Markers

- `TODO`
- `FIXME`
- `XXX`
- `WARN`
- `WARNING`
- `BUG`
- `HACK`
- `NOTE`
- `REVIEW`
- `OPTIMIZE`
- `PERF`
- `SECURITY`
- `DEPRECATED`
- `TEMP`
- `TBD`
- `TASK`

Markers can be replaced in board configuration.

## Lint Rules

- `TODO001`: a linked annotation references a missing task
- `TODO002`: a linked annotation references a done or archived task
- `TODO003`: the same task ID is linked by more than one code annotation

## Scanner Boundaries

The scanner is language-independent and intentionally heuristic. It recognizes
configured markers following common comment prefixes rather than parsing every
language grammar.

It skips configured ignored directories, nested Git checkouts, the task tree,
files larger than 2 MiB, and files with a NUL byte near the beginning. Markdown
content inside fenced code blocks is also skipped.

This approach keeps annotations portable across languages while avoiding a
language-parser dependency for every file type.
