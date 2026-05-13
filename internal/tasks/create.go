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
	State    string
	Slug     string
	Title    string
	Priority string
	Owner    string
	Tags     []string
	Now      time.Time
}

type CreateResult struct {
	Path string
	ID   string
}

func Create(repoRoot string, opts CreateOptions) (CreateResult, error) {
	cfg := DefaultConfig()
	root, err := resolveRepoRoot(repoRoot, cfg)
	if err != nil {
		return CreateResult{}, err
	}

	if !contains(cfg.States, opts.State) {
		return CreateResult{}, fmt.Errorf("unknown task state %q", opts.State)
	}
	if opts.Title == "" {
		return CreateResult{}, fmt.Errorf("task title is required")
	}
	if opts.Priority == "" {
		opts.Priority = "medium"
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

	filename := date + "-" + slug + ".md"
	path := filepath.Join(root, cfg.TaskRoot, opts.State, filename)
	if _, err := os.Stat(path); err == nil {
		return CreateResult{}, fmt.Errorf("%s already exists", filepath.ToSlash(path))
	} else if err != nil && !os.IsNotExist(err) {
		return CreateResult{}, err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return CreateResult{}, err
	}

	content := renderTask(CreateOptions{
		State:    opts.State,
		Slug:     slug,
		Title:    opts.Title,
		Priority: opts.Priority,
		Owner:    opts.Owner,
		Tags:     opts.Tags,
		Now:      opts.Now,
	}, taskID, date)

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return CreateResult{}, err
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return CreateResult{}, err
	}
	return CreateResult{Path: filepath.ToSlash(rel), ID: taskID}, nil
}

func renderTask(opts CreateOptions, taskID, date string) string {
	var tags strings.Builder
	for _, tag := range opts.Tags {
		tag = normalizeSlug(tag)
		if tag == "" {
			continue
		}
		fmt.Fprintf(&tags, "  - %s\n", tag)
	}
	if tags.Len() == 0 {
		tags.WriteString("  - task\n")
	}

	return fmt.Sprintf(`---
id: %s
title: %s
status: %s
priority: %s
owner: %s
tags:
%screated: %s
---

## Problem

Describe what needs to change, and why.

## Done when

- The expected outcome is clear
- Relevant checks pass
- The completion path is captured in Git
`, taskID, opts.Title, opts.State, opts.Priority, opts.Owner, tags.String(), date)
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
