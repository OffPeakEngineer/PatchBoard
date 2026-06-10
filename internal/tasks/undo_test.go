package tasks

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestUndoPreviewsAndRestoresTrackedTaskChanges(t *testing.T) {
	root := gitRepo(t)
	writeFile(t, root, "tasks/ready/fix-login.md", `---
id: task-auth
---

# Fix login
`)
	git(t, root, "add", ".")
	git(t, root, "-c", "user.name=Patchboard", "-c", "user.email=patchboard@example.invalid", "commit", "-m", "initial")

	if _, err := Move(root, MoveOptions{Task: "fix-login", State: "done"}); err != nil {
		t.Fatalf("Move returned error: %v", err)
	}

	preview, err := Undo(root, UndoOptions{})
	if err != nil {
		t.Fatalf("Undo preview returned error: %v", err)
	}
	if preview.Applied || len(preview.Changes) == 0 {
		t.Fatalf("expected unapplied preview with changes, got %#v", preview)
	}
	assertPathExists(t, root, "tasks", "done", "fix-login.md")

	applied, err := Undo(root, UndoOptions{Apply: true})
	if err != nil {
		t.Fatalf("Undo apply returned error: %v", err)
	}
	if !applied.Applied {
		t.Fatalf("expected undo to apply, got %#v", applied)
	}
	assertPathExists(t, root, "tasks", "ready", "fix-login.md")
	if _, err := os.Stat(filepath.Join(root, "tasks", "done", "fix-login.md")); !os.IsNotExist(err) {
		t.Fatalf("expected moved path to be restored away, got %v", err)
	}
}

func TestUndoRefusesUntrackedTaskFiles(t *testing.T) {
	root := gitRepo(t)
	writeFile(t, root, "tasks/ready/.gitkeep", "")
	git(t, root, "add", ".")
	git(t, root, "-c", "user.name=Patchboard", "-c", "user.email=patchboard@example.invalid", "commit", "-m", "initial")
	writeFile(t, root, "tasks/ready/new-task.md", "# New task\n")

	result, err := Undo(root, UndoOptions{Apply: true})
	if err != nil {
		t.Fatalf("Undo returned error: %v", err)
	}
	if !result.Refused || result.Applied {
		t.Fatalf("expected undo to refuse untracked task files, got %#v", result)
	}
	assertPathExists(t, root, "tasks", "ready", "new-task.md")
}

func gitRepo(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	git(t, root, "init")
	return root
}

func git(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(out))
	}
}
