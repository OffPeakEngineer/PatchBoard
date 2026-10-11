package tasks

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestLoadConfigThroughSymlinkUsesTargetBoardAndBoundary(t *testing.T) {
	workspace := t.TempDir()
	project := t.TempDir()
	writeFile(t, workspace, "tasks/ready/parent.md", "# Parent\n")
	writeFile(t, project, ".git", "gitdir: elsewhere\n")
	link := filepath.Join(workspace, "infrastructure")
	if err := os.Symlink(project, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, _, err := LoadConfig(link); err == nil {
		t.Fatal("linked project without a board must not select the workspace board")
	}
	writeFile(t, project, "tasks/ready/task.md", "# Task\n")
	for _, start := range []string{link, filepath.Join(link, "tasks")} {
		root, _, err := LoadConfig(start)
		if err != nil || root != project {
			t.Fatalf("LoadConfig(%q) = %q, %v; want %q", start, root, err, project)
		}
	}
}

func TestLoadConfigRejectsMissingStartingDirectory(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/ready/task.md", "# Task\n")
	if _, _, err := LoadConfig(filepath.Join(root, "infrastructures", "tasks")); err == nil {
		t.Fatal("missing explicit directory must not select an ancestor board")
	}
}

func TestLoadConfigSupportsLanePathsAndTerminalStates(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", `schema_version: 0
lanes:
  - id: anti-feature
    path: -1_anti-feature
    terminal: true
  - id: todo
    path: 1_todo
  - id: done
    path: 5_done
    terminal: true
`)
	_, cfg, err := LoadConfig(root)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(cfg.States, []string{"-1_anti-feature", "1_todo", "5_done"}) ||
		!slices.Equal(cfg.DoneStates, []string{"-1_anti-feature", "5_done"}) {
		t.Fatalf("unexpected lane configuration: %#v", cfg)
	}
}

func TestLaneConfigWithoutTerminalStatesDoesNotUseDefaultDoneStates(t *testing.T) {
	cfg := DefaultConfig()
	if err := unmarshalConfig([]byte("lanes:\n  - id: ready\n    path: 1_ready\n"), &cfg); err != nil {
		t.Fatal(err)
	}
	normalizeConfig(&cfg)
	if len(cfg.DoneStates) != 0 {
		t.Fatalf("nonterminal lanes inherited done states: %v", cfg.DoneStates)
	}
}

func TestExplicitStatesTakePrecedenceOverLanes(t *testing.T) {
	cfg := DefaultConfig()
	if err := unmarshalConfig([]byte("states: [ready, archived]\ndone_states: [archived]\nlanes:\n  - path: ignored\n    terminal: true\n"), &cfg); err != nil {
		t.Fatal(err)
	}
	normalizeConfig(&cfg)
	if !slices.Equal(cfg.States, []string{"ready", "archived"}) || !slices.Equal(cfg.DoneStates, []string{"archived"}) {
		t.Fatalf("explicit states were overridden: %#v", cfg)
	}
}

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
