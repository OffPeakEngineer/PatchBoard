package tasks

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

var BoardConfigFileNames = []string{"board.yml", "board.yaml", "board.json"}
var RootConfigFileNames = []string{".patchboard.yaml", ".patchboard.yml", ".patchboard.json"}

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
		States:     []string{"backlog", "ready", "doing", "blocked", "done", "archived"},
		DoneStates: []string{"done", "archived"},
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
	if err := unmarshalConfig(configPath, body, &cfg); err != nil {
		return "", Config{}, err
	}
	normalizeConfig(&cfg)
	return root, cfg, nil
}

func unmarshalConfig(path string, body []byte, cfg *Config) error {
	if strings.HasSuffix(path, ".json") {
		return json.Unmarshal(body, cfg)
	}
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

	for {
		for _, configFileName := range BoardConfigFileNames {
			configPath := filepath.Join(current, cfg.TaskRoot, configFileName)
			if _, err := os.Stat(configPath); err == nil {
				return current, configPath, nil
			} else if err != nil && !errors.Is(err, os.ErrNotExist) {
				return "", "", err
			}
		}

		for _, configFileName := range RootConfigFileNames {
			configPath := filepath.Join(current, configFileName)
			if _, err := os.Stat(configPath); err == nil {
				return current, configPath, nil
			} else if err != nil && !errors.Is(err, os.ErrNotExist) {
				return "", "", err
			}
		}

		info, err := os.Stat(filepath.Join(current, cfg.TaskRoot))
		if err == nil && info.IsDir() {
			return current, "", nil
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", "", err
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", "", os.ErrNotExist
		}
		current = parent
	}
}
