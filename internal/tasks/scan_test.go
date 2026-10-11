// Copyright (c) 2026 Andrew David LeTourneau
// SPDX-License-Identifier: MIT OR Zlib

package tasks

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLintSkipsSymlinksDuringTraversal(t *testing.T) {
	root := t.TempDir()
	external := t.TempDir()
	writeFile(t, root, "tasks/2_doing/task.md", "# Task\n")
	writeFile(t, root, "src/local.go", "// TODO: local follow-up\n")
	writeFile(t, external, "remote.go", "// TODO: external follow-up\n")
	for link, target := range map[string]string{
		"external":                external,
		"remote.go":               filepath.Join(external, "remote.go"),
		"broken":                  filepath.Join(external, "missing"),
		"tasks/2_doing/linked.md": external,
	} {
		if err := os.Symlink(target, filepath.Join(root, link)); err != nil {
			t.Skipf("symlinks unavailable: %v", err)
		}
	}
	result, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Tasks) != 1 || len(result.Todos) != 1 || result.Todos[0].Path != "src/local.go" {
		t.Fatalf("unexpected scan result: %#v", result)
	}
}

func TestLintDoesNotTreatTaskTemplatesAsTasks(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/2_doing/task.md", "---\nid: task-real\n---\n# Task\n")
	writeFile(t, root, "tasks/templates/example.md", "---\nid: task-real\n---\n# Template\n")
	result, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Tasks) != 1 || hasIssue(result, "TASK001") || hasIssue(result, "TASK002") {
		t.Fatalf("template was scanned as a task: %#v", result)
	}
}
