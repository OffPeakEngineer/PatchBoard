package tasks

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type FixOptions struct {
	DryRun bool
}

type FixResult struct {
	RepoRoot   string         `json:"repo_root"`
	TaskRoot   string         `json:"task_root"`
	DryRun     bool           `json:"dry_run"`
	Operations []FixOperation `json:"operations"`
}

type FixOperation struct {
	Code    string `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
	Applied bool   `json:"applied"`
}

func Fix(repoRoot string, opts FixOptions) (FixResult, error) {
	root, cfg, err := LoadConfig(repoRoot)
	if err != nil {
		return FixResult{}, err
	}

	taskRoot := filepath.Join(root, cfg.TaskRoot)
	result := FixResult{
		RepoRoot: filepath.ToSlash(root),
		TaskRoot: filepath.ToSlash(taskRoot),
		DryRun:   opts.DryRun,
	}

	if err := addScaffoldFixes(root, cfg, &result); err != nil {
		return FixResult{}, err
	}

	taskList, _, err := scanRoot(root, cfg)
	if err != nil {
		return FixResult{}, err
	}
	for _, task := range taskList {
		if task.FrontmatterStat == "" {
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(task.Path))
		operation := FixOperation{
			Code:    "FIX_STATUS",
			Path:    task.Path,
			Message: "remove redundant frontmatter status; folder location is authoritative",
			Applied: !opts.DryRun,
		}
		if !opts.DryRun {
			changed, err := removeFrontmatterStatus(path)
			if err != nil {
				return FixResult{}, err
			}
			if !changed {
				continue
			}
		}
		result.Operations = append(result.Operations, operation)
	}

	return result, nil
}

func addScaffoldFixes(root string, cfg Config, result *FixResult) error {
	taskRoot := filepath.Join(root, cfg.TaskRoot)
	if _, err := addMkdirFix(taskRoot, "FIX_TASK_ROOT", "create task root directory", result); err != nil {
		return err
	}
	for _, state := range cfg.States {
		stateDir := filepath.Join(taskRoot, state)
		created, err := addMkdirFix(stateDir, "FIX_STATE_DIR", fmt.Sprintf("create configured state directory %q", state), result)
		if err != nil {
			return err
		}
		keepFile := filepath.Join(stateDir, ".gitkeep")
		if created || dirIsEmpty(stateDir, keepFile) {
			if err := addFileFix(keepFile, "", "FIX_GITKEEP", "create missing .gitkeep", result); err != nil {
				return err
			}
		}
	}

	readmeContent, err := taskReadme(cfg)
	if err != nil {
		return err
	}
	if err := addFileFix(filepath.Join(taskRoot, "README.md"), readmeContent, "FIX_README", "install missing task board README", result); err != nil {
		return err
	}

	kanbanContent, err := templateFile("kanban.html")
	if err != nil {
		return err
	}
	return addTemplateFileFix(filepath.Join(taskRoot, "kanban.html"), kanbanContent, "FIX_KANBAN", "install or update kanban.html from template", result)
}

func addMkdirFix(path, code, message string, result *FixResult) (bool, error) {
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return false, fmt.Errorf("%s exists and is not a directory", path)
		}
		return false, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return false, err
	}

	operation := FixOperation{
		Code:    code,
		Path:    filepath.ToSlash(path),
		Message: message,
		Applied: !result.DryRun,
	}
	if !result.DryRun {
		if err := os.MkdirAll(path, 0o755); err != nil {
			return false, err
		}
	}
	result.Operations = append(result.Operations, operation)
	return true, nil
}

func addFileFix(path, content, code, message string, result *FixResult) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	operation := FixOperation{
		Code:    code,
		Path:    filepath.ToSlash(path),
		Message: message,
		Applied: !result.DryRun,
	}
	if !result.DryRun {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	result.Operations = append(result.Operations, operation)
	return nil
}

func addTemplateFileFix(path, content, code, message string, result *FixResult) error {
	actual, err := os.ReadFile(path)
	if err == nil && bytes.Equal(actual, []byte(content)) {
		return nil
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	operation := FixOperation{
		Code:    code,
		Path:    filepath.ToSlash(path),
		Message: message,
		Applied: !result.DryRun,
	}
	if !result.DryRun {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}
	result.Operations = append(result.Operations, operation)
	return nil
}

func dirIsEmpty(path, ignoredPath string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if entryPath == ignoredPath {
			continue
		}
		return false
	}
	return true
}

func removeFrontmatterStatus(path string) (bool, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	if !bytes.HasPrefix(body, []byte("---\n")) {
		return false, nil
	}
	end := bytes.Index(body[len("---\n"):], []byte("\n---"))
	if end < 0 {
		return false, nil
	}
	end += len("---\n")

	frontmatter := body[:end]
	rest := body[end:]
	pattern := regexp.MustCompile(`(?m)^status:\s*.*\n?`)
	updatedFrontmatter := pattern.ReplaceAll(frontmatter, nil)
	if bytes.Equal(frontmatter, updatedFrontmatter) {
		return false, nil
	}
	if bytes.HasSuffix(updatedFrontmatter, []byte("\n")) && bytes.HasPrefix(rest, []byte("\n---")) {
		updatedFrontmatter = bytes.TrimSuffix(updatedFrontmatter, []byte("\n"))
	}

	updated := append([]byte{}, updatedFrontmatter...)
	updated = append(updated, rest...)
	if err := os.WriteFile(path, updated, 0o644); err != nil {
		return false, err
	}
	return true, nil
}
