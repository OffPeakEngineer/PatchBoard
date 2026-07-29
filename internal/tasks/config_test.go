package tasks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigStopsAtNestedGitBoundary(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/ready/parent.md", "# Parent task\n")
	writeFile(t, root, "nested/.git", "gitdir: elsewhere\n")
	start := filepath.Join(root, "nested", "src")
	if err := os.MkdirAll(start, 0o755); err != nil {
		t.Fatal(err)
	}

	_, _, err := LoadConfig(start)
	if err == nil {
		t.Fatal("expected nested repository not to discover the parent board")
	}
	if !strings.Contains(err.Error(), filepath.Join(root, "nested")) {
		t.Fatalf("expected error to identify Git boundary, got %v", err)
	}
}

func TestLoadConfigFindsNonGitBoardFromNestedDirectory(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/ready/task.md", "# Task\n")
	start := filepath.Join(root, "src", "nested")
	if err := os.MkdirAll(start, 0o755); err != nil {
		t.Fatal(err)
	}

	found, _, err := LoadConfig(start)
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if found != root {
		t.Fatalf("found root %q, want %q", found, root)
	}
}

func TestResolveInitRootUsesGitWorktreeRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}
	start := filepath.Join(root, "src", "nested")
	if err := os.MkdirAll(start, 0o755); err != nil {
		t.Fatal(err)
	}

	found, err := ResolveInitRoot(start, false)
	if err != nil {
		t.Fatalf("ResolveInitRoot returned error: %v", err)
	}
	if found != root {
		t.Fatalf("found root %q, want %q", found, root)
	}
}

func TestResolveInitRootRequiresForceOutsideGit(t *testing.T) {
	root := t.TempDir()
	if _, err := ResolveInitRoot(root, false); err == nil {
		t.Fatal("expected non-Git initialization to require force")
	}
	found, err := ResolveInitRoot(root, true)
	if err != nil {
		t.Fatalf("ResolveInitRoot with force returned error: %v", err)
	}
	if found != root {
		t.Fatalf("found root %q, want %q", found, root)
	}
}
