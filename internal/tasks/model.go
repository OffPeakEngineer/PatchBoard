package tasks

import (
	"fmt"
	"strings"
)

// Task is the normalized representation of one Markdown file in the task
// board. The scanner keeps both derived fields and frontmatter fields so lint
// rules can explain mismatches clearly.
type Task struct {
	ID               string            `json:"id"`
	Title            string            `json:"title"`
	State            string            `json:"state"`
	Path             string            `json:"path"`
	Sources          map[string]string `json:"sources,omitempty"`
	FrontmatterID    string            `json:"frontmatter_id,omitempty"`
	FrontmatterTitle string            `json:"frontmatter_title,omitempty"`
	FrontmatterStat  string            `json:"frontmatter_status,omitempty"`
}

// Annotation is a code comment marker such as TODO, FIXME, XXX, or WARN.
// It may point at a Patchboard task with MARKER[task-id]:, or it may simply
// record loose follow-up text with MARKER:.
type Annotation struct {
	Marker string `json:"marker"`
	TaskID string `json:"task_id,omitempty"`
	Owner  string `json:"owner,omitempty"`
	Path   string `json:"path"`
	Line   int    `json:"line"`
	Syntax string `json:"syntax"`
	Text   string `json:"text,omitempty"`
}

func (t Annotation) String() string {
	location := t.Path
	if t.Line > 0 {
		location = fmt.Sprintf("%s:%d", location, t.Line)
	}

	annotation := t.Marker
	if t.TaskID != "" {
		annotation = fmt.Sprintf("%s[%s]", annotation, t.TaskID)
	}
	if t.Owner != "" {
		annotation = fmt.Sprintf("%s (@%s)", annotation, t.Owner)
	}
	if t.Text != "" {
		return fmt.Sprintf("%s  %s: %s", location, annotation, t.Text)
	}
	return fmt.Sprintf("%s  %s", location, annotation)
}

// Issue is one lint finding formatted for human output. Path is repo-relative
// when the issue came from a scanned file.
type Issue struct {
	Severity string `json:"severity"`
	Code     string `json:"code"`
	Path     string `json:"path"`
	Line     int    `json:"line,omitempty"`
	Message  string `json:"message"`
}

func (i Issue) String() string {
	location := i.Path
	if i.Line > 0 {
		location = fmt.Sprintf("%s:%d", location, i.Line)
	}
	return fmt.Sprintf("%s %s %s  %s", strings.ToUpper(i.Severity), i.Code, location, i.Message)
}

// Result is the complete scan and lint output. Commands can use Tasks and Todos
// for read-only views, and Issues for validation output.
type Result struct {
	Tasks  []Task       `json:"tasks"`
	Todos  []Annotation `json:"todos"`
	Issues []Issue      `json:"issues"`
}

// HasErrors reports whether the result should fail the CLI with exit code 1.
// Warnings can be added later without changing this decision point.
func (r Result) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}
