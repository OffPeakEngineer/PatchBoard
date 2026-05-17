package tasks

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLintDetectsTaskAndTodoIssues(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/doing/fix-login.md", `---
id: task-auth
title: Fix login
status: done
---

## Fix login
`)
	writeFile(t, root, "tasks/ready/other.md", `---
id: task-auth
title: Other
---
`)
	writeFile(t, root, "src/session.go", `package src

// `+todoRef("missing-task")+`: wire this to a real task
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}

	for _, code := range []string{"TASK002", "TASK004", "TODO001"} {
		if !hasIssue(result, code) {
			t.Fatalf("expected issue %s, got %#v", code, result.Issues)
		}
	}
}

func TestLintFlagsDoneTaskWithTodoReference(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/done/fix-login.md", `---
id: task-auth
title: Fix login
status: done
---
`)
	writeFile(t, root, "src/session.go", `package src

// `+todoRef("task-auth")+`: remove this once the redirect path is verified
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !hasIssue(result, "TODO002") {
		t.Fatalf("expected TODO002, got %#v", result.Issues)
	}
}

func TestScanFindsCaseInsensitiveAnnotationsAcrossCommentStyles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/doing/fix-login.md", `---
id: task-auth
title: Fix login
status: doing
---
`)
	writeFile(t, root, "src/session.go", `package src

// `+markerRef("warn", "task-auth")+`: this is linked to a task
# `+markerOwner("todo", "dave")+`: this is owned but not task-linked
<!-- `+markerBare("XXX")+`: this still counts as an annotation -->
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if len(result.Todos) != 3 {
		t.Fatalf("expected 3 annotations, got %#v", result.Todos)
	}
	if result.Todos[0].Marker != "WARN" || result.Todos[0].TaskID != "task-auth" {
		t.Fatalf("expected linked WARN annotation, got %#v", result.Todos[0])
	}
	if result.Todos[0].Text != "this is linked to a task" {
		t.Fatalf("expected annotation text, got %#v", result.Todos[0])
	}
	if result.Todos[1].Marker != "TODO" || result.Todos[1].Owner != "dave" || result.Todos[1].TaskID != "" {
		t.Fatalf("expected owned unlinked TODO annotation, got %#v", result.Todos[1])
	}
	if result.HasErrors() {
		t.Fatalf("expected clean lint, got %#v", result.Issues)
	}
}

func TestLintIgnoresNestedGitCheckouts(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/doing/root-task.md", `---
id: task-root
title: Root task
status: doing
---
`)
	writeFile(t, root, "nested/.git", "gitdir: ../.git/modules/nested\n")
	writeFile(t, root, "nested/src/file.go", `package src

// `+todoRef("missing-nested-task")+`: this belongs to the nested checkout
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if result.HasErrors() {
		t.Fatalf("expected nested checkout to be ignored, got %#v", result.Issues)
	}
}

func TestLintCanWarnAboutConfiguredFilenamePattern(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".patchboard.yaml", `
filename:
  enabled: true
  pattern: "^\\d{4}-\\d{2}-\\d{2}-[a-z0-9]+(?:-[a-z0-9]+)*\\.md$"
  description: YYYY-MM-DD-slug.md
  severity: warning
`)
	writeFile(t, root, "tasks/ready/nice-human-note.md", `---
id: task-filename
title: Filename note
status: ready
---
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !hasIssue(result, "TASK005") {
		t.Fatalf("expected TASK005, got %#v", result.Issues)
	}
	if result.HasErrors() {
		t.Fatalf("expected filename warning not to fail lint, got %#v", result.Issues)
	}
}

func TestLintCanLoadJsonConfig(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".patchboard.json", `{
  "filename": {
    "enabled": true,
    "pattern": "^ticket-[a-z0-9-]+\\.md$",
    "description": "ticket-slug.md",
    "severity": "warning"
  }
}
`)
	writeFile(t, root, "tasks/ready/human-note.md", `---
id: task-json-config
title: JSON config
status: ready
---
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !hasIssue(result, "TASK005") {
		t.Fatalf("expected TASK005, got %#v", result.Issues)
	}
}

func TestInitCreatesDefaultTaskBoardWithoutOverwriting(t *testing.T) {
	root := t.TempDir()

	result, err := Init(root)
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}
	if len(result.Created) == 0 {
		t.Fatal("expected init to create paths")
	}

	for _, state := range DefaultConfig().States {
		assertPathExists(t, root, "tasks", state, ".gitkeep")
	}
	assertPathExists(t, root, "tasks", "README.md")

	readme := filepath.Join(root, "tasks", "README.md")
	if err := os.WriteFile(readme, []byte("custom docs\n"), 0o644); err != nil {
		t.Fatalf("customizing README: %v", err)
	}
	second, err := Init(root)
	if err != nil {
		t.Fatalf("second Init returned error: %v", err)
	}
	if len(second.Existing) == 0 {
		t.Fatal("expected second init to report existing paths")
	}
	body, err := os.ReadFile(readme)
	if err != nil {
		t.Fatalf("reading README: %v", err)
	}
	if string(body) != "custom docs\n" {
		t.Fatalf("init overwrote README: %q", string(body))
	}
}

func TestCreateDefaultsToFirstConfiguredState(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".patchboard.yaml", `
states:
  - 0_backlog
  - 1_ready
  - 2_doing
done_states:
  - 2_doing
`)

	result, err := Create(root, CreateOptions{
		Title: "Default state task",
		Now:   time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if result.Path != "tasks/0_backlog/2026-05-17-default-state-task.md" {
		t.Fatalf("unexpected path: %s", result.Path)
	}
}

func writeFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("creating parent dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writing file: %v", err)
	}
}

func assertPathExists(t *testing.T, parts ...string) {
	t.Helper()
	path := filepath.Join(parts...)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func hasIssue(result Result, code string) bool {
	for _, issue := range result.Issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func todoRef(id string) string {
	return "TODO" + "[" + id + "]"
}

func markerRef(marker, id string) string {
	return marker + "[" + id + "]"
}

func markerOwner(marker, owner string) string {
	return marker + " (" + owner + ")"
}

func markerBare(marker string) string {
	return marker
}
