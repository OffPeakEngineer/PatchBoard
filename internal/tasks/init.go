package tasks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type InitResult struct {
	TaskRoot string
	Created  []string
	Existing []string
}

// Init creates the default task board structure under repoRoot. It is careful
// to leave existing files alone so it can be run repeatedly in a repo or added
// to bootstrap scripts without risking local task notes.
func Init(repoRoot string) (InitResult, error) {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return InitResult{}, err
	}
	cfg, err := loadConfigAtRoot(root)
	if err != nil {
		return InitResult{}, err
	} else if found && configPath != "" {
		body, err := os.ReadFile(configPath)
		if err != nil {
			return InitResult{}, err
		}
		if err := unmarshalConfig(body, &cfg); err != nil {
			return InitResult{}, err
		}
		normalizeConfig(&cfg)
	}

	result := InitResult{TaskRoot: filepath.ToSlash(filepath.Join(root, cfg.TaskRoot))}
	taskRoot := filepath.Join(root, cfg.TaskRoot)
	if err := mkdirTracked(taskRoot, &result); err != nil {
		return InitResult{}, err
	}

	for _, state := range cfg.States {
		stateDir := filepath.Join(taskRoot, state)
		if err := mkdirTracked(stateDir, &result); err != nil {
			return InitResult{}, err
		}
		keepFile := filepath.Join(stateDir, ".gitkeep")
		if err := writeFileIfMissing(keepFile, "", &result); err != nil {
			return InitResult{}, err
		}
	}

	readme := filepath.Join(taskRoot, "README.md")
	readmeContent, err := taskReadme(cfg)
	if err != nil {
		return InitResult{}, err
	}
	if err := writeFileIfMissing(readme, readmeContent, &result); err != nil {
		return InitResult{}, err
	}

	kanban := filepath.Join(taskRoot, "kanban.html")
	kanbanContent, err := templateFile("kanban.html")
	if err != nil {
		return InitResult{}, err
	}
	if err := writeFileIfMissing(kanban, kanbanContent, &result); err != nil {
		return InitResult{}, err
	}

	return result, nil
}

func mkdirTracked(path string, result *InitResult) error {
	if info, err := os.Stat(path); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%s exists and is not a directory", path)
		}
		result.Existing = append(result.Existing, filepath.ToSlash(path))
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}
	result.Created = append(result.Created, filepath.ToSlash(path))
	return nil
}

func writeFileIfMissing(path, content string, result *InitResult) error {
	if _, err := os.Stat(path); err == nil {
		result.Existing = append(result.Existing, filepath.ToSlash(path))
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	result.Created = append(result.Created, filepath.ToSlash(path))
	return nil
}

func taskReadme(cfg Config) (string, error) {
	return renderTemplate("README.md.tmpl", struct {
		States []string
	}{
		States: cfg.States,
	})
}
