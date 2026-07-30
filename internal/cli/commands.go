package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"ledoerr/patchboard/internal/tasks"
)

func runInvocation(inv cliInvocation) {
	var err error
	switch inv.command {
	case "create":
		result, err := tasks.Create(inv.repoRoot, tasks.CreateOptions{
			State: inv.state, Slug: inv.slug, Title: inv.title,
			Owner: inv.owner, Tags: splitTags(inv.tags),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("Created %s (%s)\n", result.Path, result.ID)
	case "move":
		result, err := tasks.Move(inv.repoRoot, tasks.MoveOptions{Task: inv.task, State: inv.state})
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		printMoveResult(result)
	case "start":
		_, cfg, err := tasks.LoadConfig(inv.repoRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		state := activeState(cfg.States)
		if state == "" {
			fmt.Fprintln(os.Stderr, "patchboard: no active state is configured")
			os.Exit(2)
		}
		result, err := tasks.Move(inv.repoRoot, tasks.MoveOptions{Task: inv.task, State: state})
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		printMoveResult(result)
	case "done":
		_, cfg, err := tasks.LoadConfig(inv.repoRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		state := doneState(cfg)
		if state == "" {
			fmt.Fprintln(os.Stderr, "patchboard: no done state is configured")
			os.Exit(2)
		}
		result, err := tasks.Move(inv.repoRoot, tasks.MoveOptions{Task: inv.task, State: state})
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		printMoveResult(result)
	case "init":
		root := inv.repoRoot
		if !inv.repoExplicit {
			root, err = tasks.ResolveInitRoot(root, inv.force)
			if err != nil {
				fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
				os.Exit(2)
			}
		}
		result, err := tasks.Init(root)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("Initialized task board at %s\n", result.TaskRoot)
		if len(result.Created) > 0 {
			fmt.Printf("Created %d paths\n", len(result.Created))
		}
	case "lint":
		result, err := tasks.Lint(inv.repoRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}

		if inv.json {
			if err := writeJSON(result); err != nil {
				fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
				os.Exit(2)
			}
		} else {
			for _, issue := range result.Issues {
				fmt.Println(issue.String())
			}
		}

		if result.HasErrors() {
			os.Exit(1)
		}
	case "doctor":
		if err := runDoctor(inv.repoRoot); err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
	case "fix":
		result, err := tasks.Fix(inv.repoRoot, tasks.FixOptions{DryRun: inv.dryRun})
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		printFixResult(result)
	case "undo":
		result, err := tasks.Undo(inv.repoRoot, tasks.UndoOptions{Apply: inv.apply})
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		printUndoResult(result)
		if result.Refused {
			os.Exit(1)
		}
	case "status":
		if err := runStatus(statusOptions{repoRoot: inv.repoRoot, json: inv.json, short: inv.short}); err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
	case "list":
		if err := runList(listOptions{repoRoot: inv.repoRoot, state: inv.state, json: inv.json}); err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
	case "todos":
		result, err := tasks.Lint(inv.repoRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		if inv.json {
			if err := writeJSON(struct {
				Todos []tasks.Annotation `json:"todos"`
			}{Todos: result.Todos}); err != nil {
				fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
				os.Exit(2)
			}
		} else {
			for _, todo := range result.Todos {
				fmt.Println(todo.String())
			}
		}
	default:
		fmt.Fprintf(os.Stderr, "patchboard: unknown command %q\n", inv.command)
		os.Exit(2)
	}
}

type listOptions struct {
	repoRoot string
	state    string
	json     bool
}

type statusOptions struct {
	repoRoot string
	json     bool
	short    bool
}

func runStatus(opts statusOptions) error {
	root, cfg, err := tasks.LoadConfig(opts.repoRoot)
	if err != nil {
		return err
	}
	result, err := tasks.Lint(root)
	if err != nil {
		return err
	}

	tasksByState := groupTasksByState(result.Tasks)
	errorCount, warningCount := issueCounts(result.Issues)

	if opts.json {
		return writeJSON(statusJSON(cfg, result, tasksByState, errorCount, warningCount))
	}

	fmt.Println("Patchboard status")
	if !opts.short {
		fmt.Printf("Repo: %s\n", filepath.ToSlash(root))
		fmt.Printf("Task root: %s\n", filepath.ToSlash(filepath.Join(root, cfg.TaskRoot)))
	}
	fmt.Println()
	fmt.Println("Tasks")
	for _, state := range cfg.States {
		fmt.Printf("  %-16s %d\n", state, len(tasksByState[state]))
	}

	activeState := activeState(cfg.States)
	if activeState != "" {
		fmt.Println()
		fmt.Printf("Active (%s)\n", activeState)
		activeTasks := tasksByState[activeState]
		sortTasks(activeTasks)
		if len(activeTasks) == 0 {
			fmt.Println("  No active tasks")
		} else {
			for _, task := range activeTasks {
				fmt.Printf("  - %s  %s\n", task.ID, task.Title)
			}
		}
	}

	fmt.Println()
	fmt.Printf("Annotations: %d\n", len(result.Todos))
	if len(result.Issues) == 0 {
		fmt.Println("Lint: clean")
		return nil
	}

	fmt.Printf("Lint: %d error(s), %d warning(s)\n", errorCount, warningCount)
	for _, issue := range result.Issues {
		fmt.Printf("  %s\n", issue.String())
	}
	return nil
}

func runList(opts listOptions) error {
	root, cfg, err := tasks.LoadConfig(opts.repoRoot)
	if err != nil {
		return err
	}
	if opts.state != "" && !isState(cfg.States, opts.state) {
		return fmt.Errorf("unknown state %q", opts.state)
	}

	result, err := tasks.Lint(root)
	if err != nil {
		return err
	}

	tasksByState := groupTasksByState(result.Tasks)
	if opts.json {
		return writeJSON(listJSON(cfg, tasksByState, opts.state))
	}

	for _, state := range cfg.States {
		if opts.state != "" && opts.state != state {
			continue
		}

		stateTasks := tasksByState[state]
		sortTasks(stateTasks)
		fmt.Printf("%s (%d)\n", state, len(stateTasks))
		for _, task := range stateTasks {
			fmt.Printf("  - %s  %s  %s\n", task.ID, task.Title, task.Path)
		}
	}
	return nil
}

func runDoctor(repoRoot string) error {
	result, err := tasks.Doctor(repoRoot)
	if err != nil {
		return err
	}

	fmt.Println("Patchboard doctor")
	fmt.Printf("Repo: %s\n", result.RepoRoot)
	fmt.Printf("Task root: %s\n", result.TaskRoot)
	fmt.Println()

	if len(result.Findings) == 0 {
		fmt.Println("No findings")
		return nil
	}

	for _, finding := range result.Findings {
		fmt.Printf("%s %s %s  %s\n", strings.ToUpper(finding.Severity), finding.Code, finding.Path, finding.Message)
		if finding.Action != "" {
			fmt.Printf("  Action: %s\n", finding.Action)
		}
	}
	return nil
}

func writeJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func splitTags(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			tags = append(tags, part)
		}
	}
	return tags
}

func groupTasksByState(taskList []tasks.Task) map[string][]tasks.Task {
	tasksByState := map[string][]tasks.Task{}
	for _, task := range taskList {
		tasksByState[task.State] = append(tasksByState[task.State], task)
	}
	return tasksByState
}

type stateSummary struct {
	State string       `json:"state"`
	Count int          `json:"count"`
	Tasks []tasks.Task `json:"tasks,omitempty"`
}

type statusSummary struct {
	TaskRoot     string         `json:"task_root"`
	States       []stateSummary `json:"states"`
	ActiveState  string         `json:"active_state,omitempty"`
	ActiveTasks  []tasks.Task   `json:"active_tasks"`
	Annotations  int            `json:"annotations"`
	LintErrors   int            `json:"lint_errors"`
	LintWarnings int            `json:"lint_warnings"`
	Issues       []tasks.Issue  `json:"issues"`
}

type listSummary struct {
	States []stateSummary `json:"states"`
}

func statusJSON(cfg tasks.Config, result tasks.Result, tasksByState map[string][]tasks.Task, errorCount, warningCount int) statusSummary {
	states := make([]stateSummary, 0, len(cfg.States))
	for _, state := range cfg.States {
		stateTasks := append([]tasks.Task(nil), tasksByState[state]...)
		sortTasks(stateTasks)
		states = append(states, stateSummary{
			State: state,
			Count: len(stateTasks),
		})
	}

	active := activeState(cfg.States)
	activeTasks := append([]tasks.Task(nil), tasksByState[active]...)
	if activeTasks == nil {
		activeTasks = []tasks.Task{}
	}
	sortTasks(activeTasks)

	return statusSummary{
		TaskRoot:     cfg.TaskRoot,
		States:       states,
		ActiveState:  active,
		ActiveTasks:  activeTasks,
		Annotations:  len(result.Todos),
		LintErrors:   errorCount,
		LintWarnings: warningCount,
		Issues:       result.Issues,
	}
}

func listJSON(cfg tasks.Config, tasksByState map[string][]tasks.Task, stateFilter string) listSummary {
	states := make([]stateSummary, 0, len(cfg.States))
	for _, state := range cfg.States {
		if stateFilter != "" && stateFilter != state {
			continue
		}
		stateTasks := append([]tasks.Task(nil), tasksByState[state]...)
		sortTasks(stateTasks)
		states = append(states, stateSummary{
			State: state,
			Count: len(stateTasks),
			Tasks: stateTasks,
		})
	}
	return listSummary{States: states}
}

func sortTasks(taskList []tasks.Task) {
	sort.Slice(taskList, func(i, j int) bool {
		if taskList[i].Title == taskList[j].Title {
			return taskList[i].ID < taskList[j].ID
		}
		return taskList[i].Title < taskList[j].Title
	})
}

func issueCounts(issues []tasks.Issue) (int, int) {
	var errors int
	var warnings int
	for _, issue := range issues {
		if issue.Severity == "error" {
			errors++
		} else {
			warnings++
		}
	}
	return errors, warnings
}

func printMoveResult(result tasks.MoveResult) {
	if !result.Moved {
		fmt.Printf("%s is already in %s (%s)\n", result.ID, result.To, result.OldPath)
		return
	}
	fmt.Printf("Moved %s from %s to %s\n", result.ID, result.From, result.To)
	fmt.Printf("%s -> %s\n", result.OldPath, result.NewPath)
}

func printFixResult(result tasks.FixResult) {
	if result.DryRun {
		fmt.Println("Patchboard fix preview")
	} else {
		fmt.Println("Patchboard fix")
	}
	fmt.Printf("Repo: %s\n", result.RepoRoot)
	fmt.Printf("Task root: %s\n", result.TaskRoot)
	fmt.Println()

	if len(result.Operations) == 0 {
		fmt.Println("No safe repairs")
		return
	}
	for _, operation := range result.Operations {
		verb := "would fix"
		if operation.Applied {
			verb = "fixed"
		}
		fmt.Printf("%s %s %s  %s\n", strings.ToUpper(verb), operation.Code, operation.Path, operation.Message)
	}
}

func printUndoResult(result tasks.UndoResult) {
	if result.Apply {
		fmt.Println("Patchboard undo")
	} else {
		fmt.Println("Patchboard undo preview")
	}
	fmt.Printf("Repo: %s\n", result.RepoRoot)
	fmt.Printf("Task root: %s\n", result.TaskRoot)
	fmt.Println()

	if len(result.Changes) == 0 {
		fmt.Println("No task-board changes to restore")
		return
	}
	for _, change := range result.Changes {
		fmt.Printf("%s %s\n", change.Status, change.Path)
	}
	fmt.Println()
	if result.Refused {
		fmt.Printf("Refused: %s\n", result.RefuseCause)
		return
	}
	if result.Applied {
		fmt.Println("Restored task-board changes with git.")
		return
	}
	fmt.Printf("To apply: %s\n", strings.Join(result.Command, " "))
	fmt.Println("Or run: patchboard undo --apply")
}

func activeState(states []string) string {
	for _, preferred := range []string{"2_doing", "doing"} {
		if isState(states, preferred) {
			return preferred
		}
	}
	if len(states) > 2 {
		return states[2]
	}
	if len(states) > 0 {
		return states[len(states)-1]
	}
	return ""
}

func doneState(cfg tasks.Config) string {
	for _, preferred := range []string{"3_done", "done", "archived"} {
		if isState(cfg.DoneStates, preferred) && isState(cfg.States, preferred) {
			return preferred
		}
	}
	for _, state := range cfg.DoneStates {
		if isState(cfg.States, state) {
			return state
		}
	}
	return ""
}

func isState(states []string, value string) bool {
	for _, state := range states {
		if value == state {
			return true
		}
	}
	return false
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: patchboard [status|create|move|start|done|doctor|fix|undo|init|list|lint|todos] [args]")
}
