# Task Files

A PatchBoard task is a Markdown file inside a configured state folder:

```text
tasks/doing/fix-login-timeout.md
```

The file should remain useful to someone reading it without PatchBoard.

## Recommended Shape

```markdown
---
id: task-20260512-auth-timeout
owner: andy
tags:
  - auth
  - bug
created: 2026-05-12
---

# Fix login timeout handling

## Problem

Users can get stuck after their session expires.

## Done when

- Expired sessions redirect cleanly
- Existing session refresh behavior still works
- Regression test added
```

Frontmatter is optional. Use it for durable metadata; use Markdown for the task
context collaborators need to understand.

## Identity

PatchBoard resolves task fields in this order:

- ID: frontmatter `id`, otherwise the filename slug
- Title: first H1 heading, then frontmatter `title`, then filename slug
- State: the immediate parent folder under the configured task root

New task files should prefer an H1 heading for the title.

## State Belongs in the Path

Because the folder is authoritative, frontmatter status is redundant. This is
invalid:

```text
tasks/done/fix-login-timeout.md
```

```yaml
status: doing
```

`patchboard lint` reports mismatches, and `patchboard fix` can remove legacy
status fields safely.

## Filenames

The default board does not enforce a naming pattern. A repository can configure
a filename rule as a warning or error. PatchBoard's own dogfood board uses
`slug.md`.

Filename rules should help collaborators recognize files, not turn ordinary
task editing into an implementation-specific puzzle. Prefer a clear description
such as `slug.md` or `YYYY-MM-DD-slug.md` in configuration.

## Moving Tasks

Moving the Markdown file changes its state and leaves the transition visible in
Git history. Use the CLI, the writable browser board, or an ordinary filesystem
move. The durable result is the same.
