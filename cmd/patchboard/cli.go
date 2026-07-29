package main

import (
	"fmt"
	"os"

	"github.com/integrii/flaggy"
)

type cliInvocation struct {
	command      string
	repoRoot     string
	repoExplicit bool
	state        string
	task         string
	slug         string
	title        string
	owner        string
	tags         string
	json         bool
	dryRun       bool
	apply        bool
	force        bool
	short        bool
}

func parseCLI(args []string) (cliInvocation, error) {
	parser := flaggy.NewParser("patchboard")
	parser.Description = "Repo-native task board and annotation linter"
	parser.ShowVersionWithVersionFlag = false
	parser.ShowCompletion = true

	var explicitRepo string
	parser.String(&explicitRepo, "C", "repo", "Repository or directory to operate on")

	commands := map[string]*flaggy.Subcommand{}
	add := func(name, description string) *flaggy.Subcommand {
		command := flaggy.NewSubcommand(name)
		command.Description = description
		parser.AttachSubcommand(command, 1)
		commands[name] = command
		return command
	}

	var inv cliInvocation
	var legacyRepo string

	create := add("create", "Create a task")
	create.String(&inv.state, "", "state", "Task workflow state")
	create.String(&inv.slug, "", "slug", "Task filename slug")
	create.String(&inv.title, "", "title", "Task title")
	inv.owner = "andy"
	create.String(&inv.owner, "", "owner", "Task owner")
	create.String(&inv.tags, "", "tags", "Comma-separated task tags")
	create.AddPositionalValue(&legacyRepo, "repo-root", 1, false, "Legacy repository path")

	move := add("move", "Move a task to a workflow state")
	move.AddPositionalValue(&inv.task, "task", 1, true, "Task ID, slug, filename, or path")
	move.AddPositionalValue(&inv.state, "state", 2, true, "Destination workflow state")
	move.AddPositionalValue(&legacyRepo, "repo-root", 3, false, "Legacy repository path")

	start := add("start", "Move a task to the active state")
	start.AddPositionalValue(&inv.task, "task", 1, true, "Task ID, slug, filename, or path")
	start.AddPositionalValue(&legacyRepo, "repo-root", 2, false, "Legacy repository path")

	done := add("done", "Move a task to the primary done state")
	done.AddPositionalValue(&inv.task, "task", 1, true, "Task ID, slug, filename, or path")
	done.AddPositionalValue(&legacyRepo, "repo-root", 2, false, "Legacy repository path")

	initCommand := add("init", "Initialize a task board")
	initCommand.Bool(&inv.force, "", "force", "Initialize the current directory outside Git")
	initCommand.AddPositionalValue(&legacyRepo, "repo-root", 1, false, "Legacy repository path")

	lint := add("lint", "Validate tasks and linked annotations")
	lint.Bool(&inv.json, "", "json", "Emit JSON")
	lint.AddPositionalValue(&legacyRepo, "repo-root", 1, false, "Legacy repository path")

	doctor := add("doctor", "Explain board setup and lint findings")
	doctor.AddPositionalValue(&legacyRepo, "repo-root", 1, false, "Legacy repository path")

	fix := add("fix", "Apply safe mechanical board repairs")
	fix.Bool(&inv.dryRun, "", "dry-run", "Preview repairs without changing files")
	fix.AddPositionalValue(&legacyRepo, "repo-root", 1, false, "Legacy repository path")

	undo := add("undo", "Preview or restore task-board changes with Git")
	undo.Bool(&inv.apply, "", "apply", "Apply the restore")
	undo.AddPositionalValue(&legacyRepo, "repo-root", 1, false, "Legacy repository path")

	status := add("status", "Summarize board health")
	status.Bool(&inv.json, "", "json", "Emit JSON")
	status.Bool(&inv.short, "s", "short", "Omit repository context")
	status.AddPositionalValue(&legacyRepo, "repo-root", 1, false, "Legacy repository path")

	list := add("list", "List tasks by workflow state")
	list.Bool(&inv.json, "", "json", "Emit JSON")
	list.AddPositionalValue(&inv.state, "state", 1, false, "Workflow state")
	list.AddPositionalValue(&legacyRepo, "repo-root", 2, false, "Legacy repository path")

	todos := add("todos", "List code annotations")
	todos.Bool(&inv.json, "", "json", "Emit JSON")
	todos.AddPositionalValue(&legacyRepo, "repo-root", 1, false, "Legacy repository path")

	if len(args) == 0 {
		inv.command = "status"
	} else if err := parser.ParseArgs(args); err != nil {
		return cliInvocation{}, err
	}

	if inv.command == "" {
		for name, command := range commands {
			if command.Used {
				inv.command = name
				break
			}
		}
	}
	if inv.command == "" {
		return cliInvocation{}, fmt.Errorf("no command selected")
	}

	// Preserve `patchboard list PATH`, where a single directory argument meant
	// the legacy repository rather than a state.
	if inv.command == "list" && legacyRepo == "" && inv.state != "" {
		if info, err := os.Stat(inv.state); err == nil && info.IsDir() {
			legacyRepo = inv.state
			inv.state = ""
		}
	}

	if explicitRepo != "" && legacyRepo != "" {
		return cliInvocation{}, fmt.Errorf("repository specified by both --repo/-C and legacy positional argument")
	}
	inv.repoExplicit = explicitRepo != "" || legacyRepo != ""
	switch {
	case explicitRepo != "":
		inv.repoRoot = explicitRepo
	case legacyRepo != "":
		inv.repoRoot = legacyRepo
	default:
		inv.repoRoot = "."
	}
	return inv, nil
}
