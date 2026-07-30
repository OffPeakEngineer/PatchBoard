package tasks

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
)

// Lint scans a repository with the default configuration and returns all
// discovered tasks, annotations, and lint issues. The CLI decides how to print
// this result and which exit code to use.
func Lint(repoRoot string) (Result, error) {
	root, cfg, err := LoadConfig(repoRoot)
	if err != nil {
		return Result{}, err
	}
	taskList, todoList, err := scanRoot(root, cfg)
	if err != nil {
		return Result{}, err
	}

	result := Result{
		Tasks:  nonNilTasks(taskList),
		Todos:  nonNilAnnotations(todoList),
		Issues: []Issue{},
	}
	result.Issues = append(result.Issues, lintTasks(taskList, cfg)...)
	result.Issues = append(result.Issues, lintTodos(taskList, todoList, cfg)...)
	result.Issues = append(result.Issues, lintTemplateFiles(root, cfg)...)
	return result, nil
}

func nonNilTasks(taskList []Task) []Task {
	if taskList == nil {
		return []Task{}
	}
	return taskList
}

func nonNilAnnotations(todoList []Annotation) []Annotation {
	if todoList == nil {
		return []Annotation{}
	}
	return todoList
}

// lintTasks validates the Markdown task board itself. The folder name is the
// canonical task state, so frontmatter status is invalid derived metadata.
func lintTasks(taskList []Task, cfg Config) []Issue {
	var issues []Issue
	seenIDs := map[string]Task{}

	for _, task := range taskList {
		if !slices.Contains(cfg.States, task.State) {
			issues = append(issues, Issue{
				Severity: "error",
				Code:     "TASK001",
				Path:     task.Path,
				Message:  fmt.Sprintf("unknown state folder %q", task.State),
			})
		}
		if previous, ok := seenIDs[task.ID]; ok {
			issues = append(issues, Issue{
				Severity: "error",
				Code:     "TASK002",
				Path:     task.Path,
				Message:  fmt.Sprintf("duplicate task id %q, first seen in %s", task.ID, previous.Path),
			})
		}
		seenIDs[task.ID] = task

		if task.Title == "" {
			issues = append(issues, Issue{
				Severity: "error",
				Code:     "TASK003",
				Path:     task.Path,
				Message:  "missing title",
			})
		}

		if task.FrontmatterStat != "" {
			issues = append(issues, Issue{
				Severity: "error",
				Code:     "TASK004",
				Path:     task.Path,
				Message:  "frontmatter status is derived from the containing folder; remove it",
			})
		}

		if cfg.Filename.Enabled && cfg.Filename.Pattern != "" {
			pattern, err := regexp.Compile(cfg.Filename.Pattern)
			if err != nil {
				issues = append(issues, Issue{
					Severity: "error",
					Code:     "TASK005",
					Path:     task.Path,
					Message:  fmt.Sprintf("filename rule pattern is invalid: %v", err),
				})
			} else if !pattern.MatchString(filepath.Base(task.Path)) {
				issues = append(issues, Issue{
					Severity: cfg.Filename.Severity,
					Code:     "TASK005",
					Path:     task.Path,
					Message:  fmt.Sprintf("filename should match %s", cfg.Filename.Description),
				})
			}
		}
	}

	return issues
}

// lintTodos only enforces annotations that explicitly link to a task ID. Loose
// unlinked annotations are useful inventory, but they should not make a repo
// fail lint until they claim a Patchboard task relationship.
func lintTodos(taskList []Task, todoList []Annotation, cfg Config) []Issue {
	var issues []Issue
	taskByID := map[string]Task{}
	todoByID := map[string]Annotation{}

	for _, task := range taskList {
		taskByID[task.ID] = task
	}

	for _, todo := range todoList {
		if todo.TaskID == "" {
			continue
		}

		task, ok := taskByID[todo.TaskID]
		if !ok {
			issues = append(issues, Issue{
				Severity: "error",
				Code:     "TODO001",
				Path:     todo.Path,
				Line:     todo.Line,
				Message:  fmt.Sprintf("%s references missing task id %q", todo.Syntax, todo.TaskID),
			})
		} else if slices.Contains(cfg.DoneStates, task.State) {
			issues = append(issues, Issue{
				Severity: "error",
				Code:     "TODO002",
				Path:     todo.Path,
				Line:     todo.Line,
				Message:  fmt.Sprintf("%s references %q, but the task is %s", todo.Syntax, todo.TaskID, task.State),
			})
		}

		if previous, ok := todoByID[todo.TaskID]; ok {
			issues = append(issues, Issue{
				Severity: "error",
				Code:     "TODO003",
				Path:     todo.Path,
				Line:     todo.Line,
				Message:  fmt.Sprintf("duplicate annotation task id %q, first seen in %s:%d", todo.TaskID, previous.Path, previous.Line),
			})
		}
		todoByID[todo.TaskID] = todo
	}

	return issues
}

func lintTemplateFiles(root string, cfg Config) []Issue {
	kanbanPath := filepath.Join(root, cfg.TaskRoot, "kanban.html")
	expected, err := templateFile("kanban.html")
	if err != nil {
		return []Issue{{
			Severity: "error",
			Code:     "KANBAN001",
			Path:     filepath.ToSlash(filepath.Join(cfg.TaskRoot, "kanban.html")),
			Message:  fmt.Sprintf("kanban template cannot be loaded: %v", err),
		}}
	}
	actual, err := os.ReadFile(kanbanPath)
	if os.IsNotExist(err) {
		return []Issue{{
			Severity: "warning",
			Code:     "KANBAN001",
			Path:     filepath.ToSlash(filepath.Join(cfg.TaskRoot, "kanban.html")),
			Message:  "kanban.html is missing; run patchboard fix to install it",
		}}
	}
	if err != nil {
		return []Issue{{
			Severity: "error",
			Code:     "KANBAN001",
			Path:     filepath.ToSlash(filepath.Join(cfg.TaskRoot, "kanban.html")),
			Message:  fmt.Sprintf("kanban.html cannot be read: %v", err),
		}}
	}
	if !bytes.Equal(actual, []byte(expected)) {
		return []Issue{{
			Severity: "warning",
			Code:     "KANBAN001",
			Path:     filepath.ToSlash(filepath.Join(cfg.TaskRoot, "kanban.html")),
			Message:  "kanban.html differs from the installed Patchboard template; run patchboard fix to update it",
		}}
	}
	return nil
}
