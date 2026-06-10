# Tasks

This is Patchboard's own dogfood board. Tasks are Markdown files, and state is
derived from the containing folder.

## States

- `-1_anti-feature/`
- `0_planning/`
- `1_ready/`
- `2_doing/`
- `3_done/`

Filenames use the repo-local convention configured in `tasks/board.yml`:

```text
pN[--rc-X.Y.Z]--icon--slug.md
```

The Markdown file remains the durable task record. The filename is for quick
scanning, and Git history is the audit trail.

Frontmatter should keep durable metadata such as `id`, `owner`, `tags`, and
`created`. Lane/status, title, and priority are derived from the path or
Markdown content when the board convention can do that without losing context.
