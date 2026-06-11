package tasks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLintDetectsTaskAndTodoIssues(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/2_doing/fix-login.md", `---
id: task-auth
title: Fix login
status: done
---

## Fix login
`)
	writeFile(t, root, "tasks/1_ready/other.md", `---
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
	writeFile(t, root, "tasks/3_done/fix-login.md", `---
id: task-auth
title: Fix login
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
	writeFile(t, root, "tasks/2_doing/fix-login.md", `---
id: task-auth
title: Fix login
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

func TestScanTreatsParenthesizedAnnotationValueAsOwner(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/2_doing/fix-login.md", `---
id: task-auth
title: Fix login
---
`)
	writeFile(t, root, "src/session.go", `package src

// `+"TO"+"DO"+`(task-auth): this is an owner/tag note, not a task link
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if len(result.Todos) != 1 {
		t.Fatalf("expected 1 annotation, got %#v", result.Todos)
	}
	if result.Todos[0].TaskID != "" || result.Todos[0].Owner != "task-auth" {
		t.Fatalf("expected parenthesized value to be owner only, got %#v", result.Todos[0])
	}
	if result.HasErrors() {
		t.Fatalf("expected clean lint, got %#v", result.Issues)
	}
}

func TestLintIgnoresNestedGitCheckouts(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/2_doing/root-task.md", `---
id: task-root
title: Root task
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
	writeFile(t, root, "tasks/board.yml", `
filename:
  enabled: true
  pattern: "^\\d{4}-\\d{2}-\\d{2}-[a-z0-9]+(?:-[a-z0-9]+)*\\.md$"
  description: YYYY-MM-DD-slug.md
  severity: warning
`)
	writeFile(t, root, "tasks/1_ready/nice-human-note.md", `---
id: task-filename
title: Filename note
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

func TestLintIgnoresRootConfigFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, ".patchboard.yaml", `
states:
  - wrong
filename:
  enabled: true
  pattern: "^ticket-[a-z0-9-]+\\.md$"
  severity: error
`)
	writeFile(t, root, "tasks/1_ready/human-note.md", `---
id: task-root-config
title: Root config
---
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if hasIssue(result, "TASK001") || hasIssue(result, "TASK005") {
		t.Fatalf("expected root config to be ignored, got %#v", result.Issues)
	}
}

func TestLintLoadsTaskBoardConfig(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", `
states:
  - 0_planning
  - 1_ready
  - 2_doing
done_states:
  - 2_doing
`)
	writeFile(t, root, "tasks/0_planning/task.md", `---
id: task-board-config
---

# Board config
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if hasIssue(result, "TASK001") {
		t.Fatalf("expected tasks/board.yml to configure states, got %#v", result.Issues)
	}
}

func TestLintWarnsWhenKanbanTemplateDrifts(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/1_ready/task.md", `---
id: task-board-config
---

# Board config
`)
	writeFile(t, root, "tasks/kanban.html", "<!doctype html><title>custom</title>\n")

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !hasIssue(result, "KANBAN001") {
		t.Fatalf("expected KANBAN001, got %#v", result.Issues)
	}
	if result.HasErrors() {
		t.Fatalf("expected kanban drift warning not to fail lint, got %#v", result.Issues)
	}
}

func TestLintWarnsWhenKanbanTemplateIsMissing(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/1_ready/task.md", `---
id: task-board-config
---

# Board config
`)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if !hasIssue(result, "KANBAN001") {
		t.Fatalf("expected KANBAN001, got %#v", result.Issues)
	}
}

func TestLintAcceptsInstalledKanbanTemplate(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/1_ready/task.md", `---
id: task-board-config
---

# Board config
`)
	kanban, err := templateFile("kanban.html")
	if err != nil {
		t.Fatalf("loading template: %v", err)
	}
	writeFile(t, root, "tasks/kanban.html", kanban)

	result, err := Lint(root)
	if err != nil {
		t.Fatalf("Lint returned error: %v", err)
	}
	if hasIssue(result, "KANBAN001") {
		t.Fatalf("expected installed kanban template to be clean, got %#v", result.Issues)
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
	for _, state := range []string{"backlog", "ready", "doing", "blocked", "done", "archived"} {
		assertPathMissing(t, root, "tasks", state)
	}
	assertPathExists(t, root, "tasks", "README.md")
	assertPathExists(t, root, "tasks", "kanban.html")

	readme := filepath.Join(root, "tasks", "README.md")
	kanban := filepath.Join(root, "tasks", "kanban.html")
	if err := os.WriteFile(readme, []byte("custom docs\n"), 0o644); err != nil {
		t.Fatalf("customizing README: %v", err)
	}
	if err := os.WriteFile(kanban, []byte("custom board\n"), 0o644); err != nil {
		t.Fatalf("customizing kanban: %v", err)
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
	body, err = os.ReadFile(kanban)
	if err != nil {
		t.Fatalf("reading kanban: %v", err)
	}
	if string(body) != "custom board\n" {
		t.Fatalf("init overwrote kanban: %q", string(body))
	}
}

func TestCreateDefaultsToFirstConfiguredState(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", `
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
	if result.Path != "tasks/0_backlog/default-state-task.md" {
		t.Fatalf("unexpected path: %s", result.Path)
	}
}

func TestCreateDefaultsToPlanningBeforeAntiFeature(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", `
states:
  - -1_anti-feature
  - 0_planning
  - 1_ready
  - 2_doing
done_states:
  - -1_anti-feature
`)

	result, err := Create(root, CreateOptions{
		Title: "Planning task",
		Now:   time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if result.Path != "tasks/0_planning/planning-task.md" {
		t.Fatalf("unexpected path: %s", result.Path)
	}
}

func TestCreateRendersDurableFrontmatterOnly(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/1_ready/.gitkeep", "")

	result, err := Create(root, CreateOptions{
		State: "1_ready",
		Title: "Generated task",
		Now:   time.Date(2026, 5, 17, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	body, err := os.ReadFile(filepath.Join(root, result.Path))
	if err != nil {
		t.Fatalf("reading task: %v", err)
	}
	for _, forbidden := range []string{"title:", "status:", "priority:"} {
		if strings.Contains(string(body), forbidden) {
			t.Fatalf("generated task includes derived field %q:\n%s", forbidden, string(body))
		}
	}
	if !strings.Contains(string(body), "# Generated task") {
		t.Fatalf("generated task should preserve title as Markdown heading:\n%s", string(body))
	}
}

func TestMoveTaskByIDBetweenConfiguredStates(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/1_ready/fix-login.md", `---
id: task-auth
---

# Fix login
`)

	result, err := Move(root, MoveOptions{Task: "task-auth", State: "2_doing"})
	if err != nil {
		t.Fatalf("Move returned error: %v", err)
	}
	if !result.Moved || result.From != "1_ready" || result.To != "2_doing" {
		t.Fatalf("unexpected move result: %#v", result)
	}
	assertPathExists(t, root, "tasks", "2_doing", "fix-login.md")
	if _, err := os.Stat(filepath.Join(root, "tasks", "1_ready", "fix-login.md")); !os.IsNotExist(err) {
		t.Fatalf("expected old task path to be gone, got %v", err)
	}
}

func TestMoveTaskCanMatchSlugAndRelativePath(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/1_ready/fix-login.md", `---
id: task-auth
---

# Fix login
`)

	if _, err := Move(root, MoveOptions{Task: "fix-login", State: "2_doing"}); err != nil {
		t.Fatalf("Move by slug returned error: %v", err)
	}
	if _, err := Move(root, MoveOptions{Task: "2_doing/fix-login.md", State: "3_done"}); err != nil {
		t.Fatalf("Move by relative path returned error: %v", err)
	}
	assertPathExists(t, root, "tasks", "3_done", "fix-login.md")
}

func TestMoveTaskRejectsAmbiguousQuery(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/1_ready/fix-login.md", `---
id: task-one
---

# Fix login
`)
	writeFile(t, root, "tasks/2_doing/fix-login.md", `---
id: task-two
---

# Fix login
`)

	if _, err := Move(root, MoveOptions{Task: "fix-login", State: "3_done"}); err == nil {
		t.Fatal("expected ambiguous task error")
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

func assertPathMissing(t *testing.T, parts ...string) {
	t.Helper()
	path := filepath.Join(parts...)
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected %s to be absent", path)
	} else if !os.IsNotExist(err) {
		t.Fatalf("checking %s: %v", path, err)
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
