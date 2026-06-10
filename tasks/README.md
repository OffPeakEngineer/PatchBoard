# Tasks

This is Patchboard's own dogfood board. Tasks are Markdown files, and state is
derived from the containing folder.

The board is the durable planning surface for this repository. Product context,
technical notes, and implementation follow-up should stay readable here even if
the Patchboard binary is not installed. Git history is the audit trail for how
the plan moved.

## States

- `-1_anti-feature/`
- `0_planning/`
- `1_ready/`
- `2_doing/`
- `3_done/`

Filenames use the repo-local convention configured in `tasks/board.yml`:

```text
slug.md
```

The Markdown file remains the durable task record. The filename is for quick
scanning, and Git history is the audit trail.

Frontmatter should keep durable metadata such as `id`, `owner`, `tags`, and
`created`. Lane/status is derived from the path, and title is derived from
Markdown content or the filename when the board convention can do that without
losing context.
