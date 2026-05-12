package tasks

type Config struct {
	TaskRoot          string
	States            []string
	DoneStates        []string
	AnnotationMarkers []string
	IgnoreDirs        []string
}

// DefaultConfig is intentionally usable without a config file. Patchboard's
// first contract is that a repo can opt in by adding a tasks/ directory and
// Markdown files, not by learning a setup language.
func DefaultConfig() Config {
	return Config{
		TaskRoot:   "tasks",
		States:     []string{"backlog", "ready", "doing", "blocked", "done", "archived"},
		DoneStates: []string{"done", "archived"},
		AnnotationMarkers: []string{
			"TODO",
			"FIXME",
			"XXX",
			"WARN",
			"WARNING",
			"BUG",
			"HACK",
			"NOTE",
			"REVIEW",
			"OPTIMIZE",
			"PERF",
			"SECURITY",
			"DEPRECATED",
			"TEMP",
			"TBD",
			"TASK",
		},
		IgnoreDirs: []string{".git", "node_modules", "vendor", "dist", "build"},
	}
}
