package tasks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDoctorFindsLooseRootTaskFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/note.md", `---
id: task-loose
title: Loose task
---
`)

	result, err := Doctor(root)
	if err != nil {
		t.Fatalf("Doctor returned error: %v", err)
	}
	if !hasDoctorFinding(result, "DOC003") {
		t.Fatalf("expected DOC003, got %#v", result.Findings)
	}
}

func TestDoctorReportsMissingConfiguredStateFolders(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", `
states:
  - 0_backlog
  - 1_ready
`)
	writeFile(t, root, "tasks/0_backlog/task.md", `---
id: task-backlog
title: Backlog task
---
`)

	result, err := Doctor(root)
	if err != nil {
		t.Fatalf("Doctor returned error: %v", err)
	}
	if !hasDoctorFinding(result, "DOC002") {
		t.Fatalf("expected DOC002, got %#v", result.Findings)
	}
}

func TestFixDryRunReportsSafeRepairsWithoutMutating(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", `
states:
  - 1_ready
  - 3_done
done_states:
  - 3_done
`)
	writeFile(t, root, "tasks/1_ready/fix-login.md", `---
id: task-auth
status: done
---

# Fix login
`)

	result, err := Fix(root, FixOptions{DryRun: true})
	if err != nil {
		t.Fatalf("Fix returned error: %v", err)
	}
	if !hasFixOperation(result, "FIX_STATE_DIR") || !hasFixOperation(result, "FIX_STATUS") {
		t.Fatalf("expected scaffold and status operations, got %#v", result.Operations)
	}
	if _, err := os.Stat(filepath.Join(root, "tasks", "3_done")); !os.IsNotExist(err) {
		t.Fatalf("dry run created state directory: %v", err)
	}
	body, err := os.ReadFile(filepath.Join(root, "tasks", "1_ready", "fix-login.md"))
	if err != nil {
		t.Fatalf("reading task: %v", err)
	}
	if !strings.Contains(string(body), "status: done") {
		t.Fatalf("dry run removed frontmatter status:\n%s", string(body))
	}
}

func TestFixAppliesSafeRepairs(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", `
states:
  - 1_ready
  - 3_done
done_states:
  - 3_done
`)
	writeFile(t, root, "tasks/1_ready/fix-login.md", `---
id: task-auth
status: done
---

# Fix login
`)

	result, err := Fix(root, FixOptions{})
	if err != nil {
		t.Fatalf("Fix returned error: %v", err)
	}
	if !hasFixOperation(result, "FIX_STATE_DIR") || !hasFixOperation(result, "FIX_STATUS") {
		t.Fatalf("expected scaffold and status operations, got %#v", result.Operations)
	}
	assertPathExists(t, root, "tasks", "3_done", ".gitkeep")
	assertPathExists(t, root, "tasks", "README.md")
	assertPathExists(t, root, "tasks", "kanban.html")

	body, err := os.ReadFile(filepath.Join(root, "tasks", "1_ready", "fix-login.md"))
	if err != nil {
		t.Fatalf("reading task: %v", err)
	}
	if strings.Contains(string(body), "status:") {
		t.Fatalf("expected status to be removed:\n%s", string(body))
	}
	if strings.Contains(string(body), "id: task-auth\n\n---") {
		t.Fatalf("expected frontmatter delimiter to stay tight:\n%s", string(body))
	}
}

func TestFixUpdatesDriftedKanbanTemplate(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/1_ready/.gitkeep", "")
	writeFile(t, root, "tasks/kanban.html", "<!doctype html><title>custom</title>\n")

	result, err := Fix(root, FixOptions{})
	if err != nil {
		t.Fatalf("Fix returned error: %v", err)
	}
	if !hasFixOperation(result, "FIX_KANBAN") {
		t.Fatalf("expected FIX_KANBAN, got %#v", result.Operations)
	}

	kanban, err := templateFile("kanban.html")
	if err != nil {
		t.Fatalf("loading template: %v", err)
	}
	body, err := os.ReadFile(filepath.Join(root, "tasks", "kanban.html"))
	if err != nil {
		t.Fatalf("reading kanban: %v", err)
	}
	if string(body) != kanban {
		t.Fatal("expected fix to update kanban from template")
	}
}

func hasDoctorFinding(result DoctorResult, code string) bool {
	for _, finding := range result.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}

func hasFixOperation(result FixResult, code string) bool {
	for _, operation := range result.Operations {
		if operation.Code == code {
			return true
		}
	}
	return false
}
