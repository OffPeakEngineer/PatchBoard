# PatchBoard Manifest

PatchBoard exists because a project should not be allowed to present a cleaner
planning story than its implementation actually supports.

Feature plans describe intended work. Source code contains the reality of that
work: compromises, deferred repairs, warnings, incomplete paths, and decisions
that do not fit neatly into a high-level ticket. When those records live in
separate systems, the implementation details become invisible to the people
deciding what happens next.

PatchBoard keeps them together.

## The Repository Should Tell the Whole Story

A project's tasks should be traceable in the same history as the work they
describe. Plans, implementation changes, discovered follow-up, and deliberate
decisions should remain legible together in a checkout of the repository.

The repository is not merely where completed code is delivered. It is where the
project's actual condition can be inspected.

## Markdown Is the Record

Tasks are ordinary Markdown files. They should be useful in an editor, a pull
request, a local clone, or an archive without requiring PatchBoard to interpret
them.

PatchBoard may make the files easier to create, move, validate, and view, but it
must not become the only way to understand them. If the binary disappears, the
project record remains.

## Location Is State

Workflow state belongs in the filesystem. Moving a task between folders changes
its state and leaves that transition visible in Git history.

PatchBoard avoids redundant status metadata because duplicated truth eventually
disagrees. The path is authoritative because it is already the action a
collaborator performed.

## Git Is Memory

Git already records changes, authorship, review, movement, and time. PatchBoard
uses that history rather than inventing another audit trail.

Planning changes deserve the same durability as implementation changes. A task
that was created, revised, rejected, completed, or moved should remain a normal
history question.

## Implementation Reality Is Planning Context

A `TODO`, `FIXME`, warning, or other code annotation is evidence about the
project's condition. It should not remain hidden from the people responsible for
planning simply because it was discovered at implementation depth.

Linked annotations connect an exact place in the code to a durable task record.
Linting keeps that relationship honest: referenced work must exist, completed
work should not leave stale promises behind, and implementation follow-up should
not silently fall outside the board.

This is the feature that makes PatchBoard more than a file-based kanban. The
board is the approachable surface; lint prevents the implementation truth from
being excluded.

## The Tool Serves the Files

The CLI and browser board are replaceable interfaces over the same durable task
tree.

Engineering collaborators should be able to use the CLI, shell, editor, Git,
and code review. Product and project collaborators should be able to use a
familiar visual board. Neither interface owns the data, and neither should make
the other second-class.

## Opinion Should Reduce Ceremony

PatchBoard is opinionated where an opinion makes the record clearer:

- folders define workflow state;
- Markdown carries durable context;
- Git carries history;
- explicit task links connect code annotations to planned work;
- lint observes and reports without silently creating work;
- mechanical repair is separated from decisions that require human judgment.

Configuration exists so a repository can express its own workflow without
abandoning those principles.

## Rejected Work Is Still Knowledge

A durable project record includes decisions not to build something. An
anti-feature can preserve reasoning, prevent repeated debate, and explain why
an apparently obvious feature is intentionally absent.

PatchBoard values legible decisions, not merely a growing list of completed
items.

## The Boundary Is Deliberate

PatchBoard is the repository-local layer of truth. It does not need to become a
portfolio manager, customer-support system, organization-wide reporting
database, or hosted collaboration platform.

Its value comes from being small enough to remain close to the work and durable
enough to outlive any particular interface. Features should strengthen
traceability, legibility, or the connection between plans and implementation.
Features that merely reproduce a general-purpose tracker should be treated with
care.

## A Project Conscience

PatchBoard should make it harder for a project to forget what its code already
knows.

The board communicates intent. Git preserves memory. Lint confronts the plan
with implementation reality. Together they provide a project conscience: a
small, inspectable system that keeps declared progress accountable to the work
it represents.
