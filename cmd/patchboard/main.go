package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"ledoerr/patchboard/internal/tasks"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "create":
		createFlags := flag.NewFlagSet("create", flag.ExitOnError)
		state := createFlags.String("state", "backlog", "task workflow state")
		slug := createFlags.String("slug", "", "task filename slug")
		title := createFlags.String("title", "", "task title")
		priority := createFlags.String("priority", "medium", "task priority")
		owner := createFlags.String("owner", "andy", "task owner")
		tags := createFlags.String("tags", "", "comma-separated task tags")
		if err := createFlags.Parse(os.Args[2:]); err != nil {
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
		root := repoRootArg()
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
		result, err := tasks.Lint(repoRootArg())
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
	case "todos":
		result, err := tasks.Lint(repoRootArg())
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
		fmt.Fprintf(os.Stderr, "patchboard: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func repoRootArg() string {
	if len(os.Args) > 2 {
		return os.Args[2]
	}
	return "."
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: patchboard <init|lint|todos> [repo-root]")
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

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: patchboard <create|init|lint|todos> [repo-root]")
}
