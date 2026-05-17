# Tasks

This is Patchboard's own dogfood board. Tasks are Markdown files, and state is
derived from the containing folder.

## States

- `-1_anti-feature/`
- `0_backlog/`
- `1_ready/`
- `2_doing/`
- `3_done/`

Filenames use the repo-local convention configured in `.patchboard.yaml`:

```text
pN--icon--slug.md
```

The Markdown file remains the durable task record. The filename is for quick
scanning, and Git history is the audit trail.
