package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/off-peak.engineer/utilities/patchboard/internal/tasks"
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

func TestSummaryCLIExportsRequestedFormat(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "tasks", "2_doing"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "tasks", "2_doing", "example.md"), []byte("# Example\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, format := range []string{"yaml", "json", "csv", "html"} {
		output := filepath.Join(root, "out", "board."+format)
		result := runCLIProcess(t, "summary", "--type="+format, "-C", root, "--output", output)
		if result.code != 0 || result.stderr != "" || !strings.Contains(result.stdout, "Exported 1 tasks to ") {
			t.Fatalf("summary %s failed: %#v", format, result)
		}
		if info, err := os.Stat(output); err != nil || info.Size() == 0 {
			t.Fatalf("summary output missing: %v", err)
		}
	}
	result := runCLIProcess(t, "summary", "-C", root)
	if result.code != 0 {
		t.Fatalf("default summary failed: %#v", result)
	}
	if _, err := os.Stat(filepath.Join(root, "summary.yaml")); err != nil {
		t.Fatal(err)
	}
	invalid := runCLIProcess(t, "summary", "--type=pdf", "-C", root)
	if invalid.code != 2 || !strings.Contains(invalid.stderr, "unsupported summary type") {
		t.Fatalf("unexpected invalid-format result: %#v", invalid)
	}
}

func TestDoneStatePrefersNumberedDoneLaneOverAntiFeature(t *testing.T) {
	cfg := tasks.Config{
		States:     []string{"-1_anti-feature", "1_todo", "2_doing", "5_done"},
		DoneStates: []string{"-1_anti-feature", "5_done"},
	}
	if state := doneState(cfg); state != "5_done" {
		t.Fatalf("done state = %q, want 5_done", state)
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
	version := runCLIProcess(t, "--version")
	if version.code != 0 || version.stderr != "" || version.stdout != "Version: "+Version+"\n" {
		t.Fatalf("unexpected version result: %#v", version)
	}

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
