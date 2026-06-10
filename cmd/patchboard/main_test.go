package main

import (
	"testing"

	"ledoerr/patchboard/internal/tasks"
)

func TestParseListArgsDefaultsToCurrentRepo(t *testing.T) {
	opts, err := parseListArgs(nil)
	if err != nil {
		t.Fatalf("parseListArgs returned error: %v", err)
	}
	if opts.repoRoot != "." || opts.state != "" {
		t.Fatalf("unexpected options: %#v", opts)
	}
}

func TestParseListArgsAcceptsStateAndRepoRoot(t *testing.T) {
	root := t.TempDir()
	opts, err := parseListArgs([]string{"1_ready", root})
	if err != nil {
		t.Fatalf("parseListArgs returned error: %v", err)
	}
	if opts.repoRoot != root || opts.state != "1_ready" {
		t.Fatalf("unexpected options: %#v", opts)
	}
}

func TestParseListArgsAcceptsRepoRootWithoutState(t *testing.T) {
	root := t.TempDir()
	opts, err := parseListArgs([]string{root})
	if err != nil {
		t.Fatalf("parseListArgs returned error: %v", err)
	}
	if opts.repoRoot != root || opts.state != "" {
		t.Fatalf("unexpected options: %#v", opts)
	}
}

func TestParseListArgsAcceptsJSONFlag(t *testing.T) {
	root := t.TempDir()
	opts, err := parseListArgs([]string{"--json", "2_doing", root})
	if err != nil {
		t.Fatalf("parseListArgs returned error: %v", err)
	}
	if opts.repoRoot != root || opts.state != "2_doing" || !opts.json {
		t.Fatalf("unexpected options: %#v", opts)
	}
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

func testConfig(states, doneStates []string) tasks.Config {
	return tasks.Config{States: states, DoneStates: doneStates}
}
