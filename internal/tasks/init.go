package tasks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type InitResult struct {
	TaskRoot string
	Created  []string
	Existing []string
}

// Init creates the default task board structure under repoRoot. It is careful
// to leave existing files alone so it can be run repeatedly in a repo or added
// to bootstrap scripts without risking local task notes.
func Init(repoRoot string) (InitResult, error) {
	cfg := DefaultConfig()
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return InitResult{}, err
	}

	result := InitResult{TaskRoot: filepath.ToSlash(filepath.Join(root, cfg.TaskRoot))}
	taskRoot := filepath.Join(root, cfg.TaskRoot)
	if err := mkdirTracked(taskRoot, &result); err != nil {
		return InitResult{}, err
	}

	for _, state := range cfg.States {
		stateDir := filepath.Join(taskRoot, state)
		if err := mkdirTracked(stateDir, &result); err != nil {
			return InitResult{}, err
		}
		keepFile := filepath.Join(stateDir, ".gitkeep")
		if err := writeFileIfMissing(keepFile, "", &result); err != nil {
			return InitResult{}, err
		}
	}

	readme := filepath.Join(taskRoot, "README.md")
	if err := writeFileIfMissing(readme, taskReadme(), &result); err != nil {
		return InitResult{}, err
	}

	return result, nil
}

func mkdirTracked(path string, result *InitResult) error {
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%s exists and is not a directory", path)
		}
		result.Existing = append(result.Existing, filepath.ToSlash(path))
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}
	result.Created = append(result.Created, filepath.ToSlash(path))
	return nil
}

func writeFileIfMissing(path, content string, result *InitResult) error {
	if _, err := os.Stat(path); err == nil {
		result.Existing = append(result.Existing, filepath.ToSlash(path))
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	result.Created = append(result.Created, filepath.ToSlash(path))
	return nil
}

func taskReadme() string {
	cfg := DefaultConfig()
	var states strings.Builder
	for _, state := range cfg.States {
		fmt.Fprintf(&states, "- `%s/`\n", state)
	}

	return fmt.Sprintf(`# Tasks

This folder is a Patchboard task board. Tasks are Markdown files, and the
folder containing a task is its workflow state.

## States

%s
Move a task file between folders to change its state. Git history is the audit
trail.

## Task Shape

~~~markdown
---
id: task-YYYYMMDD-short-name
title: Short, concrete task title
status: backlog
priority: medium
owner: your-name
created: YYYY-MM-DD
---

## Problem

What needs to change, and why?

## Done when

- The expected behavior is implemented
- Relevant tests or checks pass
~~~

The folder is authoritative for status. If frontmatter includes `+"`status`"+`,
it should match the parent folder.

## Code Annotations

Link code comments back to tasks with square brackets:

~~~text
TODO[task-YYYYMMDD-short-name]: describe the follow-up
FIXME[task-YYYYMMDD-short-name]: describe the known problem
~~~

Unlinked annotations such as `+"`TODO:`"+`, `+"`XXX:`"+`, and `+"`WARN:`"+` are useful inventory,
but they do not fail lint until they reference a task ID.
`, states.String())
}
