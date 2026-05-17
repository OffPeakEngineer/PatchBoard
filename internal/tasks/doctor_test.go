package tasks

import "testing"

func TestDoctorFindsLooseRootTaskFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/p2--note.md", `---
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
	writeFile(t, root, ".patchboard.yaml", `
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

func hasDoctorFinding(result DoctorResult, code string) bool {
	for _, finding := range result.Findings {
		if finding.Code == code {
			return true
		}
	}
	return false
}
