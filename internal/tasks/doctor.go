package tasks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DoctorFinding struct {
	Severity string
	Code     string
	Path     string
	Message  string
	Action   string
}

type DoctorResult struct {
	RepoRoot string
	TaskRoot string
	Findings []DoctorFinding
}

func Doctor(repoRoot string) (DoctorResult, error) {
	root, cfg, err := LoadConfig(repoRoot)
	if err != nil {
		return DoctorResult{}, err
	}

	taskRoot := filepath.Join(root, cfg.TaskRoot)
	result := DoctorResult{
		RepoRoot: filepath.ToSlash(root),
		TaskRoot: filepath.ToSlash(filepath.Join(root, cfg.TaskRoot)),
	}

	for _, state := range cfg.States {
		statePath := filepath.Join(taskRoot, state)
		if info, err := os.Stat(statePath); err == nil {
			if !info.IsDir() {
				result.Findings = append(result.Findings, DoctorFinding{
					Severity: "error",
					Code:     "DOC001",
					Path:     filepath.ToSlash(statePath),
					Message:  fmt.Sprintf("configured state %q exists but is not a directory", state),
					Action:   "rename the file or choose a different state name",
				})
			}
		} else if errors.Is(err, os.ErrNotExist) {
			result.Findings = append(result.Findings, DoctorFinding{
				Severity: "warning",
				Code:     "DOC002",
				Path:     filepath.ToSlash(statePath),
				Message:  fmt.Sprintf("configured state %q is missing", state),
				Action:   "run patchboard init to create missing state folders",
			})
		} else {
			return DoctorResult{}, err
		}
	}

	looseTasks, err := looseTaskFiles(taskRoot, cfg.TaskRoot)
	if err != nil {
		return DoctorResult{}, err
	}
	for _, path := range looseTasks {
		result.Findings = append(result.Findings, DoctorFinding{
			Severity: "warning",
			Code:     "DOC003",
			Path:     path,
			Message:  "task file is directly under the task root, so it is not assigned to a workflow state",
			Action:   fmt.Sprintf("move it under one of: %s", strings.Join(cfg.States, ", ")),
		})
	}

	lintResult, err := Lint(root)
	if err != nil {
		return DoctorResult{}, err
	}
	for _, issue := range lintResult.Issues {
		result.Findings = append(result.Findings, DoctorFinding{
			Severity: issue.Severity,
			Code:     issue.Code,
			Path:     issue.Path,
			Message:  issue.Message,
			Action:   doctorAction(issue),
		})
	}

	return result, nil
}

func looseTaskFiles(taskRoot, taskRootName string) ([]string, error) {
	entries, err := os.ReadDir(taskRoot)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".md" || entry.Name() == "README.md" {
			continue
		}
		paths = append(paths, filepath.ToSlash(filepath.Join(taskRootName, entry.Name())))
	}
	return paths, nil
}

func doctorAction(issue Issue) string {
	switch issue.Code {
	case "TASK001":
		return "move the task into a configured state folder or add the state to config"
	case "TASK002":
		return "give one of the tasks a distinct frontmatter id"
	case "TASK003":
		return "add frontmatter title or a Markdown heading"
	case "TASK004":
		return "update frontmatter status to match the containing folder, or remove status"
	case "TASK005":
		return "rename the file or adjust the filename rule"
	case "TODO001":
		return "create the referenced task or update the annotation task id"
	case "TODO002":
		return "remove the annotation or move the referenced task out of a done state"
	case "TODO003":
		return "keep only one code annotation per task id"
	default:
		return "inspect the finding and update the task board"
	}
}
