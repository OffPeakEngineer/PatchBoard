package cli

import (
	"fmt"
	"os"
)

// Main runs the PatchBoard command-line application.
func Main() {
	inv, err := parseCLI(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "patchboard: %v\n", err)
		os.Exit(2)
	}
	runInvocation(inv)
}

// main is retained for the subprocess-based CLI tests in this package.
func main() {
	Main()
}
