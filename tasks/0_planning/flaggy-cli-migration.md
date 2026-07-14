---
id: task-20260714-flaggy-cli-migration
owner: andy
tags:
  - patchboard
  - cli
  - flaggy
  - refactor
created: 2026-07-14
---

# Migrate command-line parsing to Flaggy

## Problem

The CLI currently creates a separate standard-library flag set for each command
and maintains manual parsing and usage behavior in a growing `main.go`. This
repeats argument plumbing and makes it harder to introduce a consistent global
repository option.

Flaggy supports subcommands, positional values, global and command-local flags,
generated help, and flags in flexible positions without requiring a particular
package layout.

## Context

This is a parser and command-organization refactor, not permission to change the
existing command contract. Root-selection behavior should follow the separately
reviewed repository resolution contract.

The dependency uses the Unlicense. Its license and dependency footprint should
be recorded alongside the decision because Patchboard has an intentionally
cautious licensing policy.

## Done when

- Flaggy is added as a direct dependency after its license is accepted
- Every existing command and flag is represented in the generated command tree
- Command-specific parsing and execution are separated from the top-level entry
  point so `main.go` no longer contains the complete CLI implementation
- Existing positional task, state, and optional repository arguments remain
  compatible unless a reviewed migration explicitly changes them
- Human help output covers all commands, flags, and positional values
- Unknown commands, missing required values, and invalid flags retain clear
  errors and intentional exit codes
- CLI tests cover argument ordering, help, output streams, and exit behavior

## Depends on

- `task-20260714-repository-resolution-contract`
