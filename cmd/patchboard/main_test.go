package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"ledoerr/patchboard/internal/tasks"
)

func TestParseCLIDefaultsToStatusInCurrentRepo(t *testing.T) {
	opts, err := parseCLI(nil)
	if err != nil {
		t.Fatalf("parseCLI returned error: %v", err)
	}
	if opts.command != "status" || opts.repoRoot != "." {
		t.Fatalf("unexpected options: %#v", opts)
	}
}

func TestParseListArgsAcceptsStateAndRepoRoot(t *testing.T) {
	root := t.TempDir()
	opts, err := parseCLI([]string{"list", "1_ready", root})
	if err != nil {
		t.Fatalf("parseCLI returned error: %v", err)
	}
	if opts.repoRoot != root || opts.state != "1_ready" {
		t.Fatalf("unexpected options: %#v", opts)
	}
}

func TestParseListArgsAcceptsRepoRootWithoutState(t *testing.T) {
	root := t.TempDir()
	opts, err := parseCLI([]string{"list", root})
	if err != nil {
		t.Fatalf("parseCLI returned error: %v", err)
	}
	if opts.repoRoot != root || opts.state != "" {
		t.Fatalf("unexpected options: %#v", opts)
	}
}

func TestParseListArgsAcceptsJSONFlag(t *testing.T) {
	root := t.TempDir()
	opts, err := parseCLI([]string{"list", "--json", "2_doing", root})
	if err != nil {
		t.Fatalf("parseCLI returned error: %v", err)
	}
	if opts.repoRoot != root || opts.state != "2_doing" || !opts.json {
		t.Fatalf("unexpected options: %#v", opts)
	}
}

func TestParseCLIAcceptsGlobalRepoBeforeOrAfterCommand(t *testing.T) {
	for _, args := range [][]string{
		{"-C", "/tmp/project", "lint"},
		{"lint", "--repo", "/tmp/project"},
	} {
		opts, err := parseCLI(args)
		if err != nil {
			t.Fatalf("parseCLI(%q): %v", args, err)
		}
		if opts.command != "lint" || opts.repoRoot != "/tmp/project" || !opts.repoExplicit {
			t.Fatalf("parseCLI(%q) returned %#v", args, opts)
		}
	}
}

func TestParseCLIRejectsConflictingRepositoryArguments(t *testing.T) {
	_, err := parseCLI([]string{"lint", "legacy", "--repo", "explicit"})
	if err == nil || !strings.Contains(err.Error(), "both") {
		t.Fatalf("expected conflicting repository error, got %v", err)
	}
}

func TestParseCLIPreservesLegacyCommandShapes(t *testing.T) {
	tests := []struct {
		args    []string
		command string
		task    string
		state   string
		repo    string
	}{
		{[]string{"status", "repo"}, "status", "", "", "repo"},
		{[]string{"move", "task", "doing", "repo"}, "move", "task", "doing", "repo"},
		{[]string{"start", "task", "repo"}, "start", "task", "", "repo"},
		{[]string{"done", "task", "repo"}, "done", "task", "", "repo"},
		{[]string{"init", "repo"}, "init", "", "", "repo"},
		{[]string{"doctor", "repo"}, "doctor", "", "", "repo"},
		{[]string{"fix", "--dry-run", "repo"}, "fix", "", "", "repo"},
		{[]string{"undo", "--apply", "repo"}, "undo", "", "", "repo"},
		{[]string{"lint", "--json", "repo"}, "lint", "", "", "repo"},
		{[]string{"todos", "--json", "repo"}, "todos", "", "", "repo"},
	}
	for _, test := range tests {
		inv, err := parseCLI(test.args)
		if err != nil {
			t.Fatalf("parseCLI(%q): %v", test.args, err)
		}
		if inv.command != test.command || inv.task != test.task || inv.state != test.state || inv.repoRoot != test.repo {
			t.Fatalf("parseCLI(%q) returned %#v", test.args, inv)
		}
	}
}

func TestCLIHelpAndErrorsUseIntentionalStreamsAndExitCodes(t *testing.T) {
	help := runCLIProcess(t, "--help")
	if help.code != 0 || help.stdout != "" || !strings.Contains(help.stderr, "Subcommands:") {
		t.Fatalf("unexpected help result: %#v", help)
	}

	unknown := runCLIProcess(t, "not-a-command")
	if unknown.code != 2 || !strings.Contains(unknown.stderr, "No subcommand") || !strings.Contains(unknown.stdout, "Available subcommands") {
		t.Fatalf("unexpected unknown-command result: %#v", unknown)
	}

	missing := runCLIProcess(t, "move")
	if missing.code != 2 || !strings.Contains(missing.stderr, "Required positional") {
		t.Fatalf("unexpected missing-argument result: %#v", missing)
	}
}

type cliProcessResult struct {
	code   int
	stdout string
	stderr string
}

func runCLIProcess(t *testing.T, args ...string) cliProcessResult {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=TestCLIHelperProcess")
	cmd.Env = append(os.Environ(), "PATCHBOARD_TEST_ARGS="+strings.Join(args, "\x1f"))
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	code := 0
	if exitErr, ok := err.(*exec.ExitError); ok {
		code = exitErr.ExitCode()
	} else if err != nil {
		t.Fatalf("running helper: %v", err)
	}
	return cliProcessResult{code: code, stdout: stdout.String(), stderr: stderr.String()}
}

func TestCLIHelperProcess(t *testing.T) {
	value, ok := os.LookupEnv("PATCHBOARD_TEST_ARGS")
	if !ok {
		return
	}
	os.Args = append([]string{"patchboard"}, strings.Split(value, "\x1f")...)
	main()
}

func TestActiveStatePrefersConfiguredDoingLane(t *testing.T) {
	state := activeState([]string{"0_backlog", "1_ready", "2_doing", "3_done"})
	if state != "2_doing" {
		t.Fatalf("unexpected active state: %q", state)
	}
}

func TestDoneStatePrefersConfiguredDoneLane(t *testing.T) {
	state := doneState(testConfig(
		[]string{"-1_anti-feature", "0_planning", "1_ready", "2_doing", "3_done"},
		[]string{"-1_anti-feature", "3_done"},
	))
	if state != "3_done" {
		t.Fatalf("unexpected done state: %q", state)
	}
}

func TestStatusJSONUsesEmptyActiveTaskList(t *testing.T) {
	cfg := testConfig([]string{"0_planning", "1_ready", "2_doing", "3_done"}, []string{"3_done"})
	status := statusJSON(cfg, tasks.Result{}, map[string][]tasks.Task{}, 0, 0)
	if status.ActiveTasks == nil {
		t.Fatal("expected active_tasks to encode as an empty array, not null")
	}
}

func testConfig(states, doneStates []string) tasks.Config {
	return tasks.Config{States: states, DoneStates: doneStates}
}
