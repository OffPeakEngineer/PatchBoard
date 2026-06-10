package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type MoveOptions struct {
	Task  string
	State string
}

type MoveResult struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	OldPath string `json:"old_path"`
	NewPath string `json:"new_path"`
	Moved   bool   `json:"moved"`
}

func Move(repoRoot string, opts MoveOptions) (MoveResult, error) {
	root, cfg, err := LoadConfig(repoRoot)
	if err != nil {
		return MoveResult{}, err
	}
	if strings.TrimSpace(opts.Task) == "" {
		return MoveResult{}, fmt.Errorf("task is required")
	}
	if strings.TrimSpace(opts.State) == "" {
		return MoveResult{}, fmt.Errorf("destination state is required")
	}
	if !contains(cfg.States, opts.State) {
		return MoveResult{}, fmt.Errorf("unknown task state %q", opts.State)
	}

	taskList, _, err := scanRoot(root, cfg)
	if err != nil {
		return MoveResult{}, err
	}
	task, err := findMoveTask(taskList, opts.Task)
	if err != nil {
		return MoveResult{}, err
	}

	result := MoveResult{
		ID:      task.ID,
		From:    task.State,
		To:      opts.State,
		OldPath: task.Path,
		NewPath: filepath.ToSlash(filepath.Join(cfg.TaskRoot, opts.State, filepath.Base(task.Path))),
	}
	if task.State == opts.State {
		return result, nil
	}

	oldPath := filepath.Join(root, filepath.FromSlash(task.Path))
	newPath := filepath.Join(root, filepath.FromSlash(result.NewPath))
	if _, err := os.Stat(newPath); err == nil {
		return MoveResult{}, fmt.Errorf("%s already exists", filepath.ToSlash(newPath))
	} else if err != nil && !os.IsNotExist(err) {
		return MoveResult{}, err
	}
	if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return MoveResult{}, err
	}
	if err := os.Rename(oldPath, newPath); err != nil {
		return MoveResult{}, err
	}

	result.Moved = true
	return result, nil
}

func findMoveTask(taskList []Task, query string) (Task, error) {
	query = strings.TrimSpace(filepath.ToSlash(query))
	var matches []Task
	for _, task := range taskList {
		if taskMatches(task, query) {
			matches = append(matches, task)
		}
	}
	if len(matches) == 0 {
		return Task{}, fmt.Errorf("task %q not found", query)
	}
	if len(matches) > 1 {
		var ids []string
		for _, match := range matches {
			ids = append(ids, fmt.Sprintf("%s (%s)", match.ID, match.Path))
		}
		return Task{}, fmt.Errorf("task %q is ambiguous: %s", query, strings.Join(ids, ", "))
	}
	return matches[0], nil
}

func taskMatches(task Task, query string) bool {
	path := task.Path
	base := filepath.Base(path)
	slug := strings.TrimSuffix(base, filepath.Ext(base))
	return query == task.ID ||
		query == task.Path ||
		query == taskPathWithoutRoot(task.Path) ||
		query == base ||
		query == slug
}

func taskPathWithoutRoot(path string) string {
	_, rest, ok := strings.Cut(path, "/")
	if !ok {
		return path
	}
	return rest
}
