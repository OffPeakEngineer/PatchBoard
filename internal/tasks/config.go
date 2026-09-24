package tasks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

const BoardConfigFileName = "board.yml"

type Config struct {
	TaskRoot          string         `json:"task_root"`
	States            []string       `json:"states"`
	DoneStates        []string       `json:"done_states"`
	AnnotationMarkers []string       `json:"annotation_markers"`
	IgnoreDirs        []string       `json:"ignore_dirs"`
	Filename          FilenameConfig `json:"filename"`
}

type FilenameConfig struct {
	Enabled     bool   `json:"enabled"`
	Pattern     string `json:"pattern"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// DefaultConfig is intentionally usable without a config file. Patchboard's
// first contract is that a repo can opt in by adding a tasks/ directory and
// Markdown files, not by learning a setup language.
func DefaultConfig() Config {
	return Config{
		TaskRoot:   "tasks",
		States:     []string{"-1_anti-feature", "0_planning", "1_ready", "2_doing", "3_done"},
		DoneStates: []string{"-1_anti-feature", "3_done"},
		AnnotationMarkers: []string{
			"TODO",
			"FIXME",
			"XXX",
			"WARN",
			"WARNING",
			"BUG",
			"HACK",
			"NOTE",
			"REVIEW",
			"OPTIMIZE",
			"PERF",
			"SECURITY",
			"DEPRECATED",
			"TEMP",
			"TBD",
			"TASK",
		},
		IgnoreDirs: []string{".git", "node_modules", "vendor", "dist", "build"},
	}
}

func LoadConfig(start string) (string, Config, error) {
	cfg := DefaultConfig()
	root, configPath, err := findRepoRoot(start, cfg)
	if err != nil {
		return "", Config{}, err
	}
	if configPath == "" {
		return root, cfg, nil
	}

	body, err := os.ReadFile(configPath)
	if err != nil {
		return "", Config{}, err
	}
	if err := unmarshalConfig(body, &cfg); err != nil {
		return "", Config{}, err
	}
	normalizeConfig(&cfg)
	return root, cfg, nil
}

// loadConfigAtRoot loads configuration only from root. Unlike LoadConfig, it
// never searches parent directories, which is important when initializing a
// nested project that should own a separate task board.
func loadConfigAtRoot(root string) (Config, error) {
	cfg := DefaultConfig()
	configPath := filepath.Join(root, cfg.TaskRoot, BoardConfigFileName)
	body, err := os.ReadFile(configPath)
	if errors.Is(err, os.ErrNotExist) {
		return cfg, nil
	}
	if err != nil {
		return Config{}, err
	}
	if err := unmarshalConfig(body, &cfg); err != nil {
		return Config{}, err
	}
	normalizeConfig(&cfg)
	return cfg, nil
}

func unmarshalConfig(body []byte, cfg *Config) error {
	return yaml.Unmarshal(body, cfg)
}

func normalizeConfig(cfg *Config) {
	defaults := DefaultConfig()
	if cfg.TaskRoot == "" {
		cfg.TaskRoot = defaults.TaskRoot
	}
	if len(cfg.States) == 0 {
		cfg.States = defaults.States
	}
	if len(cfg.DoneStates) == 0 {
		cfg.DoneStates = defaults.DoneStates
	}
	if len(cfg.AnnotationMarkers) == 0 {
		cfg.AnnotationMarkers = defaults.AnnotationMarkers
	}
	if len(cfg.IgnoreDirs) == 0 {
		cfg.IgnoreDirs = defaults.IgnoreDirs
	}
	if cfg.Filename.Severity == "" {
		cfg.Filename.Severity = "warning"
	}
	if cfg.Filename.Description == "" {
		cfg.Filename.Description = cfg.Filename.Pattern
	}
}

func findRepoRoot(start string, cfg Config) (string, string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", "", err
	}
	start = current
	gitRoot, err := findGitRoot(current)
	if err != nil {
		return "", "", err
	}

	for {
		if configPath, found, err := boardAtRoot(current, cfg); err != nil {
			return "", "", err
		} else if found {
			return current, configPath, nil
		}

		if gitRoot != "" && current == gitRoot {
			break
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	boundary := "no Git worktree was found"
	if gitRoot != "" {
		boundary = fmt.Sprintf("Git worktree boundary is %s", gitRoot)
	}
	return "", "", fmt.Errorf("no Patchboard board found from %s (%s); searched for %s or %s; run patchboard init or select a repository with -C PATH", start, boundary, filepath.Join(cfg.TaskRoot, "board.yml"), cfg.TaskRoot)
}

func boardAtRoot(root string, cfg Config) (string, bool, error) {
	configPath := filepath.Join(root, cfg.TaskRoot, BoardConfigFileName)
	if _, err := os.Stat(configPath); err == nil {
		return configPath, true, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", false, err
	}
	info, err := os.Stat(filepath.Join(root, cfg.TaskRoot))
	if err == nil && info.IsDir() {
		return "", true, nil
	}
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", false, err
	}
	return "", false, nil
}

func findGitRoot(start string) (string, error) {
	current := start
	for {
		if _, err := os.Stat(filepath.Join(current, ".git")); err == nil {
			return current, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", nil
		}
		current = parent
	}
}

// ResolveInitRoot chooses the default initialization target. Existing boards
// remain valid outside Git; new non-Git boards require an explicit force.
func ResolveInitRoot(start string, force bool) (string, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	if root, _, err := findRepoRoot(abs, DefaultConfig()); err == nil {
		return root, nil
	}
	gitRoot, err := findGitRoot(abs)
	if err != nil {
		return "", err
	}
	if gitRoot != "" {
		return gitRoot, nil
	}
	if force {
		return abs, nil
	}
	return "", fmt.Errorf("cannot infer a project root from %s: no Git worktree or existing tasks directory; rerun with patchboard init --force or select a repository with -C PATH", abs)
}
