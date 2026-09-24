# Core Concepts

PatchBoard is built around a deliberately small contract: tasks are Markdown
files, workflow state comes from folders, and Git history records how the plan
changed.

## The Filesystem Is the Durable Record

PatchBoard does not hide project work in a database. A typical board looks like:

```text
tasks/
  README.md
  kanban.html
  backlog/
  ready/
  doing/
  blocked/
  done/
  archived/
```

Moving a task is a file move. Editing a task is a normal text edit. Reviewing
how a plan changed is a normal Git history question.

The CLI and browser board are helpers around that tree. If either disappears,
the task record remains usable.

## Folder State Is Authoritative

The immediate folder under the task root defines a task's state:

```text
tasks/doing/fix-login-timeout.md
```

In this example, the task is `doing`. Storing another status value in
frontmatter would duplicate information and could disagree with the path.
PatchBoard therefore treats frontmatter status as legacy data and can remove it
with `patchboard fix`.

## Planning and Implementation Share Context

High-level feature work and low-level implementation follow-up often live in
different systems. PatchBoard links them with annotations such as:

```text
TODO[task-20260512-auth-timeout]: handle expired refresh token
```

The task file carries the planning context. The annotation identifies the exact
place where work remains. Linting verifies that the relationship stays valid as
the board and code change.

## Git Is the Audit Trail

Task creation, editing, and movement appear in the same history as the work they
describe. That keeps planning changes visible in reviews and preserves them in
ordinary clones without a separate export or service.

PatchBoard's `undo` command uses Git rather than maintaining a second history
mechanism.

## Opinionated but Configurable

PatchBoard works without configuration and provides a conventional default
board. Repositories can replace the default states, done states, annotation
markers, ignored directories, and filename conventions in `tasks/board.yml`.

Configuration changes local conventions without changing the underlying file
contract.

## Product Boundary

PatchBoard is intended to be the repository-local layer of truth. It is not a
replacement for portfolio planning, customer-support queues, or organization-
wide reporting systems. Those systems can describe broad commitments while
PatchBoard records the concrete state of work in a particular repository.
