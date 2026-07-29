package main

import (
	"fmt"
	"os"
)

func main() {
	inv, err := parseCLI(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
		os.Exit(2)
	}
	runInvocation(inv)
}
