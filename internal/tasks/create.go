package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type CreateOptions struct {
	State string
	Slug  string
	Title string
	Owner string
	Tags  []string
	Now   time.Time
}

type CreateResult struct {
	Path string
	ID   string
}

func Create(repoRoot string, opts CreateOptions) (CreateResult, error) {
	root, cfg, err := LoadConfig(repoRoot)
	if err != nil {
		return CreateResult{}, err
	}

	if opts.State == "" {
		opts.State = defaultCreateState(cfg.States)
	}
	if !contains(cfg.States, opts.State) {
		return CreateResult{}, fmt.Errorf("unknown task state %q", opts.State)
	}
	if opts.Title == "" {
		return CreateResult{}, fmt.Errorf("task title is required")
	}
	if opts.Owner == "" {
		opts.Owner = "andy"
	}
	if opts.Now.IsZero() {
		opts.Now = time.Now()
	}

	slug := normalizeSlug(opts.Slug)
	if slug == "" {
		slug = normalizeSlug(opts.Title)
	}
	date := opts.Now.Format("2006-01-02")
	idDate := opts.Now.Format("20060102")
	taskID := "task-" + idDate + "-" + slug

	filename := slug + ".md"
	path := filepath.Join(root, cfg.TaskRoot, opts.State, filename)
	if _, err := os.Stat(path); err == nil {
		return CreateResult{}, fmt.Errorf("%s already exists", filepath.ToSlash(path))
	} else if err != nil && !os.IsNotExist(err) {
		return CreateResult{}, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return CreateResult{}, err
	}

	content, err := renderTask(CreateOptions{
		State: opts.State,
		Slug:  slug,
		Title: opts.Title,
		Owner: opts.Owner,
		Tags:  opts.Tags,
		Now:   opts.Now,
	}, taskID, date)
	if err != nil {
		return CreateResult{}, err
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return CreateResult{}, err
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return CreateResult{}, err
	}
	return CreateResult{Path: filepath.ToSlash(rel), ID: taskID}, nil
}

func renderTask(opts CreateOptions, taskID, date string) (string, error) {
	var tags []string
	for _, tag := range opts.Tags {
		tag = normalizeSlug(tag)
		if tag == "" {
			continue
		}
		tags = append(tags, tag)
	}
	if len(tags) == 0 {
		tags = append(tags, "task")
	}

	return renderTemplate("task.md.tmpl", struct {
		ID      string
		Owner   string
		Tags    []string
		Created string
		Title   string
	}{
		ID:      taskID,
		Owner:   opts.Owner,
		Tags:    tags,
		Created: date,
		Title:   opts.Title,
	})
}

func normalizeSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	value = regexp.MustCompile(`-+`).ReplaceAllString(value, "-")
	return value
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func defaultCreateState(states []string) string {
	for _, preferred := range []string{"0_planning", "backlog", "planning", "0_backlog", "1_ready", "ready"} {
		if contains(states, preferred) {
			return preferred
		}
	}
	if len(states) > 0 {
		return states[0]
	}
	return ""
}
