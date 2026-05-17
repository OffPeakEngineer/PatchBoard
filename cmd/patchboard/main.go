package main

import (
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
		if err := runStatus("."); err != nil {
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
		priority := createFlags.String("priority", "medium", "task priority")
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
			State:    *state,
			Slug:     *slug,
			Title:    *title,
			Priority: *priority,
			Owner:    *owner,
			Tags:     splitTags(*tags),
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
		result, err := tasks.Lint(repoRootArg(args[1:]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}

		for _, issue := range result.Issues {
			fmt.Println(issue.String())
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
		if err := runStatus(repoRootArg(args[1:])); err != nil {
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
		result, err := tasks.Lint(repoRootArg(args[1:]))
		if err != nil {
			fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
			os.Exit(2)
		}
		for _, todo := range result.Todos {
			fmt.Println(todo.String())
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
}

func runStatus(repoRoot string) error {
	root, cfg, err := tasks.LoadConfig(repoRoot)
	if err != nil {
		return err
	}
	result, err := tasks.Lint(root)
	if err != nil {
		return err
	}

	tasksByState := groupTasksByState(result.Tasks)
	errorCount, warningCount := issueCounts(result.Issues)

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
