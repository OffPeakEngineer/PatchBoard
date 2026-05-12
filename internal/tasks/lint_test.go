package tasks

import (
	"os"
	"path/filepath"
	"testing"
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
