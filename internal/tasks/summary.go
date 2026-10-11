// Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
package tasks

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"sigs.k8s.io/yaml"
)

type SummaryOptions struct {
	Type   string
	Output string
}

type SummaryResult struct {
	Path  string
	Tasks int
}

// BoardSummary is a portable, ordered snapshot of the task board. YAML and
// JSON share this schema; CSV flattens its tasks into one row per task.
type BoardSummary struct {
	TaskRoot   string         `json:"task_root"`
	TotalTasks int            `json:"total_tasks"`
	States     []SummaryState `json:"states"`
}

type SummaryState struct {
	State string        `json:"state"`
	Count int           `json:"count"`
	Tasks []SummaryTask `json:"tasks"`
}

type SummaryTask struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
	Path  string `json:"path"`
	// Raw Markdown is included only in the offline HTML document.
	Markdown string `json:"-"`
	Anchor   string `json:"-"`
	Filename string `json:"-"`
}

func ExportSummary(repoRoot string, opts SummaryOptions) (SummaryResult, error) {
	format := strings.ToLower(opts.Type)
	if format == "" {
		format = "yaml"
	}
	switch format {
	case "yaml", "json", "csv", "html":
	default:
		return SummaryResult{}, fmt.Errorf("unsupported summary type %q; choose yaml, json, csv, or html", opts.Type)
	}
	root, cfg, err := LoadConfig(repoRoot)
	if err != nil {
		return SummaryResult{}, err
	}
	board, err := summarizeBoard(root, cfg, format == "html")
	if err != nil {
		return SummaryResult{}, err
	}
	body, err := encodeSummary(board, format)
	if err != nil {
		return SummaryResult{}, err
	}
	output := opts.Output
	if output == "" {
		output = filepath.Join(root, "summary."+format)
	}
	output, err = filepath.Abs(output)
	if err != nil {
		return SummaryResult{}, err
	}
	// Refuse to replace task Markdown or editable board assets with an export.
	rel, err := filepath.Rel(filepath.Join(root, cfg.TaskRoot), output)
	if err != nil {
		return SummaryResult{}, err
	}
	insideTasks := rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
	if insideTasks && (filepath.Ext(output) == ".md" || rel == BoardConfigFileName || rel == "kanban.html") {
		return SummaryResult{}, fmt.Errorf("summary output would replace a task or editable board file: %s", output)
	}
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return SummaryResult{}, err
	}
	// A temporary file and rename keep a previous snapshot intact on failure.
	file, err := os.CreateTemp(filepath.Dir(output), ".patchboard-summary-*")
	if err != nil {
		return SummaryResult{}, err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(body); err != nil {
		file.Close()
		return SummaryResult{}, err
	}
	if err := file.Chmod(0o644); err != nil {
		file.Close()
		return SummaryResult{}, err
	}
	if err := file.Close(); err != nil {
		return SummaryResult{}, err
	}
	if err := os.Rename(file.Name(), output); err != nil {
		return SummaryResult{}, err
	}
	return SummaryResult{Path: filepath.ToSlash(output), Tasks: board.TotalTasks}, nil
}

func summarizeBoard(root string, cfg Config, includeMarkdown bool) (BoardSummary, error) {
	taskRoot := filepath.Join(root, cfg.TaskRoot)
	valid := false
	for _, state := range cfg.States {
		if info, err := os.Stat(filepath.Join(taskRoot, state)); err == nil && info.IsDir() {
			valid = true
			break
		}
	}
	if !valid {
		return BoardSummary{}, fmt.Errorf("%s is not a PatchBoard task root: no configured state directories exist", taskRoot)
	}
	taskList, err := scanTasks(taskRoot, cfg)
	if err != nil {
		return BoardSummary{}, err
	}
	sort.Slice(taskList, func(i, j int) bool {
		if taskList[i].Title != taskList[j].Title {
			return taskList[i].Title < taskList[j].Title
		}
		return taskList[i].Path < taskList[j].Path
	})
	states := append([]string(nil), cfg.States...)
	byState := map[string][]SummaryTask{}
	unknown := map[string]bool{}
	for i, task := range taskList {
		entry := SummaryTask{ID: task.ID, Title: task.Title, State: task.State, Path: task.Path,
			Anchor: "task-" + strconv.Itoa(i), Filename: filepath.Base(task.Path)}
		if includeMarkdown {
			body, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(task.Path)))
			if err != nil {
				return BoardSummary{}, err
			}
			entry.Markdown = string(body)
		}
		byState[task.State] = append(byState[task.State], entry)
		if !contains(cfg.States, task.State) {
			unknown[task.State] = true
		}
	}
	// Keep tasks in unconfigured folders visible so an export cannot silently
	// hide work while the board's configuration is being corrected.
	var extra []string
	for state := range unknown {
		extra = append(extra, state)
	}
	sort.Strings(extra)
	states = append(states, extra...)
	board := BoardSummary{TaskRoot: cfg.TaskRoot, TotalTasks: len(taskList), States: []SummaryState{}}
	for _, state := range states {
		entries := byState[state]
		if entries == nil {
			entries = []SummaryTask{}
		}
		board.States = append(board.States, SummaryState{State: state, Count: len(entries), Tasks: entries})
	}
	return board, nil
}

func encodeSummary(board BoardSummary, format string) ([]byte, error) {
	switch format {
	case "yaml":
		return yaml.Marshal(board)
	case "json":
		body, err := json.MarshalIndent(board, "", "  ")
		return append(body, '\n'), err
	case "csv":
		var body bytes.Buffer
		writer := csv.NewWriter(&body)
		if err := writer.Write([]string{"id", "title", "state", "path"}); err != nil {
			return nil, err
		}
		for _, state := range board.States {
			for _, task := range state.Tasks {
				if err := writer.Write([]string{task.ID, task.Title, task.State, task.Path}); err != nil {
					return nil, err
				}
			}
		}
		writer.Flush()
		return body.Bytes(), writer.Error()
	case "html":
		return renderSummaryHTML(board)
	}
	return nil, fmt.Errorf("unsupported summary type %q", format)
}

func renderSummaryHTML(board BoardSummary) ([]byte, error) {
	source, err := templateFile("summary.html.tmpl")
	if err != nil {
		return nil, err
	}
	kanban, err := templateFile("kanban.html")
	if err != nil {
		return nil, err
	}
	// Reuse the embedded browser board's styles, without its interactive script.
	_, styles, found := strings.Cut(kanban, "<style>")
	if !found {
		return nil, fmt.Errorf("kanban template is missing its styles")
	}
	styles, _, found = strings.Cut(styles, "</style>")
	if !found {
		return nil, fmt.Errorf("kanban template has an unclosed style block")
	}
	tmpl, err := template.New("summary").Parse(source)
	if err != nil {
		return nil, err
	}
	var body bytes.Buffer
	data := struct {
		BoardSummary
		Styles template.CSS
	}{board, template.CSS(styles)} // Styles are trusted, embedded source only.
	if err := tmpl.Execute(&body, data); err != nil {
		return nil, err
	}
	return body.Bytes(), nil
}
