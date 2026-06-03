---
id: task-20260513-patchboard-dogfood-overhaul
title: Patchboard dogfood overhaul
owner: andy
tags:
  - patchboard
  - dogfood
  - tasks
  - process
created: 2026-05-13
---

## Problem

Patchboard is useful as a repo-native paper trail, but the current dogfood pass
showed places where the tool can fight natural file-first workflows. For now,
the paper trail matters more than perfect lint conformance.

## Lessons

- Moving files between lanes should make status updates easy, forgiving, or
  automatically repairable.
  - The git history should be treated as the first-line-of-defense for tracking movement of tasks, as that will map the developers actions to code commit changes (for software projects, a naturalized citizen)
- Lint should distinguish between project-manager-friendly guidance and hard failures that block work.
- Templates may need relaxed filename rules or a dedicated template schema.
- The tool should support humans who edit files directly with Vim or an IDE.
- Patchboard may need a command that files its own follow-up tasks when it finds
  non-blocking issues.
- The lane/config model should stay KISS and work from YAML or JSON config.
- Status, priority, and title should derive from path conventions instead of
  being duplicated in frontmatter.
- Repos should be able to define numbered lanes such as `0_backlog`,
  `1_ready`, and terminal lanes such as `x_done` or `-1_closed` using a simple human template, like mustache, not regex
- Filename lint should support conventions like `p<priority>--<emoji>--<slug>.md` and optional release-candidate segments such as `p1--rc-1.0.0--🚀--launch.md`.

## Open questions

- Should status derive from folder location, frontmatter, or a reconciled
  combination of both?
  - Folder location always, unless the option is an undo.
  - The linter is expected to keep the frontmatter updated. The simple human action is to rename the title, and/or drag and drop the file. The linter should be able to see that move and complete the steps missed (update frontmatter) and not scold the user for doing it "there way" if we reliably got to the result.
- Should `patchboard lint` offer a `--fix` mode for safe frontmatter repairs?
  - No, as long as git is in a clean state, the 'repairs' can be achieved with a git checkout or a wrapper `patchboard undo` that does the git action to restore.
- Which errors should fail CI versus only produce warnings?
  - Only errors that cannot be elegantly reconciled need to be elevated, and even then not scary, more like "did you mean?". We will assume un-reconcilable issues were just uncompleted thoughts/work before someone was trying to "checkpoint", but they didn't give enough context for the "team" (aka doctor rules) cannot parse or deal with.
- How should moved files preserve history without forcing noisy edits?
  - git preserves history, not our tool
- Should Patchboard read `tasks/board.yml`, `.patchboard.yml`, or both?
  - I think it should live in  `tasks` and probably be `.patchboard.yml`, as we'll want to try and hide it from 'regular' users, the people that don't know that a `.git` folder exists too, but it's developer safe. I also want to support . `patchboard.json` (and env or TOML, if it's free. I want the team to choose the config, the tool does not need an opinion on this as long as it isn't doing extra legwork to support it)
- How should the tool report derived fields in JSON without pretending they
  were hand-authored metadata?
  - the config file will contain a way to define those fields and their names

## Done when

- Patchboard's dogfood requirements are written in its own repo.
- The infrastructure repo can keep using `/tasks` as a paper trail without
  blocking cluster work.
- Lint behavior is friendlier to direct file editing and PM workflows.
