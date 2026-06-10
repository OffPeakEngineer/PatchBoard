package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"ledoerr/patchboard/internal/tasks"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		if err := runStatus(statusOptions{repoRoot: "."}); err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		return
	}

	switch args[0] {
	case "create":
		createFlags := flag.NewFlagSet("create", flag.ExitOnError)
		state := createFlags.String("state", "", "task workflow state")
		slug := createFlags.String("slug", "", "task filename slug")
		title := createFlags.String("title", "", "task title")
		owner := createFlags.String("owner", "andy", "task owner")
		tags := createFlags.String("tags", "", "comma-separated task tags")
		if err := createFlags.Parse(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		root := "."
		if createFlags.NArg() > 0 {
			root = createFlags.Arg(0)
		}
		result, err := tasks.Create(root, tasks.CreateOptions{
			State: *state,
			Slug:  *slug,
			Title: *title,
			Owner: *owner,
			Tags:  splitTags(*tags),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		fmt.Printf("Created %s (%s)\n", result.Path, result.ID)
	case "init":
		root := repoRootArg(args[1:])
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
		opts, err := parseJSONRepoArgs("lint", args[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		result, err := tasks.Lint(opts.repoRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}

		if opts.json {
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
		if err := runDoctor(repoRootArg(args[1:])); err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
	case "status":
		opts, err := parseJSONRepoArgs("status", args[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		if err := runStatus(statusOptions{repoRoot: opts.repoRoot, json: opts.json}); err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
	case "list":
		opts, err := parseListArgs(args[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		if err := runList(opts); err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
	case "todos":
		opts, err := parseJSONRepoArgs("todos", args[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		result, err := tasks.Lint(opts.repoRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		if opts.json {
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
	case "help", "-h", "--help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "patchboard: unknown command %q\n\n", args[0])
		usage()
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
}

type jsonRepoOptions struct {
	repoRoot string
	json     bool
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

func parseListArgs(args []string) (listOptions, error) {
	opts := listOptions{repoRoot: "."}
	var filtered []string
	for _, arg := range args {
		if arg == "--json" {
			opts.json = true
			continue
		}
		filtered = append(filtered, arg)
	}
	args = filtered
	if len(args) == 0 {
		return opts, nil
	}
	if len(args) > 2 {
		return opts, fmt.Errorf("usage: patchboard list [state] [repo-root]")
	}

	if len(args) == 1 {
		if info, err := os.Stat(args[0]); err == nil && info.IsDir() {
			opts.repoRoot = args[0]
		} else {
			opts.state = args[0]
		}
		return opts, nil
	}

	opts.state = args[0]
	opts.repoRoot = args[1]
	return opts, nil
}

func parseJSONRepoArgs(name string, args []string) (jsonRepoOptions, error) {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	jsonOut := flags.Bool("json", false, "emit JSON")
	if err := flags.Parse(args); err != nil {
		return jsonRepoOptions{}, err
	}
	if flags.NArg() > 1 {
		return jsonRepoOptions{}, fmt.Errorf("usage: patchboard %s [--json] [repo-root]", name)
	}
	opts := jsonRepoOptions{repoRoot: ".", json: *jsonOut}
	if flags.NArg() == 1 {
		opts.repoRoot = flags.Arg(0)
	}
	return opts, nil
}

func writeJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func repoRootArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return "."
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

func isState(states []string, value string) bool {
	for _, state := range states {
		if value == state {
			return true
		}
	}
	return false
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: patchboard [status|create|doctor|init|list|lint|todos] [args]")
}
