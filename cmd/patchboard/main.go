package main

import (
	"fmt"
	"os"

	"ledoerr/patchboard/internal/tasks"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "init":
		root := "."
		if len(os.Args) > 2 {
			root = os.Args[2]
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
		result, err := tasks.Lint(".")
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
		result, err := tasks.Lint(".")
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

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: patchboard <init [repo-root]|lint|todos>")
}
