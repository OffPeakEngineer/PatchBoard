// Package tasks contains Patchboard's core filesystem model.
//
// The package is deliberately small: it scans Markdown task files, scans source
// comments for annotations, and returns lint issues. Command-line formatting,
// future JSON output, and any web UI should stay at the edges of the program so
// this package remains easy to test.
package tasks
