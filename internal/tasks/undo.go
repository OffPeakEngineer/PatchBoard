package tasks

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type UndoOptions struct {
	Apply bool
}

type UndoResult struct {
	RepoRoot    string      `json:"repo_root"`
	TaskRoot    string      `json:"task_root"`
	Apply       bool        `json:"apply"`
	Changes     []GitChange `json:"changes"`
	Command     []string    `json:"command"`
	Applied     bool        `json:"applied"`
	Refused     bool        `json:"refused"`
	RefuseCause string      `json:"refuse_cause,omitempty"`
}

type GitChange struct {
	Status string `json:"status"`
	Path   string `json:"path"`
}

func Undo(repoRoot string, opts UndoOptions) (UndoResult, error) {
	root, cfg, err := LoadConfig(repoRoot)
	if err != nil {
		return UndoResult{}, err
	}
	result := UndoResult{
		RepoRoot: filepath.ToSlash(root),
		TaskRoot: filepath.ToSlash(filepath.Join(root, cfg.TaskRoot)),
		Apply:    opts.Apply,
		Command:  []string{"git", "-C", root, "restore", "--source=HEAD", "--staged", "--worktree", "--", cfg.TaskRoot},
	}

	changes, err := gitStatus(root, cfg.TaskRoot)
	if err != nil {
		return UndoResult{}, err
	}
	result.Changes = changes
	if len(changes) == 0 {
		return result, nil
	}
	safeUntracked := safeUntrackedMovePaths(changes)
	for _, change := range changes {
		if change.Status == "??" && !safeUntracked[change.Path] {
			result.Refused = true
			result.RefuseCause = "untracked task files are present; inspect them before using git clean manually"
			return result, nil
		}
	}
	if !opts.Apply {
		return result, nil
	}
	for path := range safeUntracked {
		if err := runGit(root, "clean", "-f", "--", path); err != nil {
			return UndoResult{}, err
		}
	}
	if err := runGit(root, "restore", "--source=HEAD", "--staged", "--worktree", "--", cfg.TaskRoot); err != nil {
		return UndoResult{}, err
	}
	result.Applied = true
	return result, nil
}

func gitStatus(root, taskRoot string) ([]GitChange, error) {
	out, err := gitOutput(root, "status", "--porcelain=v1", "-uall", "-z", "--", taskRoot)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, nil
	}

	parts := bytes.Split(bytes.TrimSuffix(out, []byte{0}), []byte{0})
	changes := make([]GitChange, 0, len(parts))
	for i := 0; i < len(parts); i++ {
		entry := string(parts[i])
		if len(entry) < 4 {
			continue
		}
		status := entry[:2]
		path := entry[3:]
		if strings.HasPrefix(status, "R") || strings.HasPrefix(status, "C") {
			if i+1 < len(parts) {
				path = path + " -> " + string(parts[i+1])
				i++
			}
		}
		changes = append(changes, GitChange{Status: status, Path: filepath.ToSlash(path)})
	}
	return changes, nil
}

func safeUntrackedMovePaths(changes []GitChange) map[string]bool {
	deletedBasenames := map[string]bool{}
	for _, change := range changes {
		if strings.Contains(change.Status, "D") {
			deletedBasenames[filepath.Base(change.Path)] = true
		}
	}

	safe := map[string]bool{}
	for _, change := range changes {
		if change.Status == "??" && deletedBasenames[filepath.Base(change.Path)] {
			safe[change.Path] = true
		}
	}
	return safe
}

func gitOutput(root string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git %s: %v\n%s", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	return out, nil
}

func runGit(root string, args ...string) error {
	_, err := gitOutput(root, args...)
	return err
}
