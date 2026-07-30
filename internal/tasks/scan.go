package tasks

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const maxScannedFileSize = 2 * 1024 * 1024

// Scan reads the task board and code annotations from repoRoot. If repoRoot is
// inside a project, Scan walks upward until it finds the configured task root;
// this keeps `go run ./cmd/patchboard lint` working from inside this module.
func Scan(repoRoot string, cfg Config) ([]Task, []Annotation, error) {
	repoRoot, err := resolveRepoRoot(repoRoot, cfg)
	if err != nil {
		return nil, nil, err
	}
	return scanRoot(repoRoot, cfg)
}

func scanRoot(repoRoot string, cfg Config) ([]Task, []Annotation, error) {
	taskRoot := filepath.Join(repoRoot, cfg.TaskRoot)

	taskList, err := scanTasks(taskRoot, cfg)
	if err != nil {
		return nil, nil, err
	}

	todoList, err := scanTodos(repoRoot, cfg)
	if err != nil {
		return nil, nil, err
	}

	return taskList, todoList, nil
}

func scanTasks(taskRoot string, cfg Config) ([]Task, error) {
	var taskList []Task

	err := filepath.WalkDir(taskRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == taskRoot {
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}

		rel, err := filepath.Rel(taskRoot, path)
		if err != nil {
			return err
		}
		parts := strings.Split(filepath.ToSlash(rel), "/")
		if len(parts) < 2 {
			return nil
		}

		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		meta, markdown := parseMarkdown(body)
		title, titleSource := firstNonEmptyWithSource(
			[]sourcedValue{
				{Value: firstHeading(markdown), Source: "markdown"},
				{Value: meta["title"], Source: "frontmatter"},
				{Value: filenameTitle(path), Source: "path"},
			},
		)
		id, idSource := firstNonEmptyWithSource(
			[]sourcedValue{
				{Value: meta["id"], Source: "frontmatter"},
				{Value: filenameSlug(path), Source: "path"},
			},
		)

		taskList = append(taskList, Task{
			ID:               id,
			Title:            title,
			State:            parts[0],
			Path:             filepath.ToSlash(filepath.Join(cfg.TaskRoot, rel)),
			Sources:          map[string]string{"id": idSource, "title": titleSource, "state": "path"},
			FrontmatterID:    meta["id"],
			FrontmatterTitle: meta["title"],
			FrontmatterStat:  meta["status"],
		})
		return nil
	})
	if errors.Is(err, os.ErrNotExist) {
		return taskList, nil
	}
	return taskList, err
}

func scanTodos(repoRoot string, cfg Config) ([]Annotation, error) {
	pattern, err := annotationPattern(cfg.AnnotationMarkers)
	if err != nil {
		return nil, err
	}

	var refs []Annotation
	err = filepath.WalkDir(repoRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := entry.Name()
		if entry.IsDir() {
			rel, err := filepath.Rel(repoRoot, path)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			if rel == "." {
				return nil
			}
			if isNestedRepo(path, rel) {
				return filepath.SkipDir
			}
			if slices.Contains(cfg.IgnoreDirs, name) || rel == cfg.TaskRoot {
				return filepath.SkipDir
			}
			return nil
		}

		fileRefs, err := scanFileTodos(path, pattern)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return err
		}
		for i := range fileRefs {
			fileRefs[i].Path = filepath.ToSlash(rel)
		}
		refs = append(refs, fileRefs...)
		return nil
	})
	return refs, err
}

func isNestedRepo(path, rel string) bool {
	if rel == "." {
		return false
	}
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		return true
	}
	return false
}

func scanFileTodos(path string, pattern *regexp.Regexp) ([]Annotation, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if info.Size() > maxScannedFileSize {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	prefix := make([]byte, 4096)
	n, err := file.Read(prefix)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if bytes.Contains(prefix[:n], []byte{0}) {
		return nil, nil
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	var refs []Annotation
	scanner := bufio.NewScanner(file)
	line := 0
	inMarkdownFence := false
	isMarkdown := filepath.Ext(path) == ".md" || filepath.Ext(path) == ".markdown"
	for scanner.Scan() {
		line++
		text := scanner.Text()
		if isMarkdown && isFenceLine(text) {
			inMarkdownFence = !inMarkdownFence
			continue
		}
		if inMarkdownFence {
			continue
		}
		for _, ref := range parseAnnotations(text, line, path, pattern) {
			refs = append(refs, ref)
		}
	}
	return refs, scanner.Err()
}

// annotationPattern recognizes the stable part of the convention:
//
//	marker with task ID, marker with owner, or marker with plain text
//
// Comment-prefix filtering happens in parseAnnotations so this pattern can stay
// focused on the annotation grammar itself.
func annotationPattern(markers []string) (*regexp.Regexp, error) {
	escaped := make([]string, 0, len(markers))
	for _, marker := range markers {
		escaped = append(escaped, regexp.QuoteMeta(marker))
	}
	pattern := `(?i)\b(` + strings.Join(escaped, "|") + `)\s*(?:\[([^\]]+)\]|\(([^\)]+)\))?\s*:`
	return regexp.Compile(pattern)
}

func parseAnnotations(line string, lineNumber int, path string, pattern *regexp.Regexp) []Annotation {
	var refs []Annotation
	matches := pattern.FindAllStringSubmatch(line, -1)
	indexes := pattern.FindAllStringSubmatchIndex(line, -1)
	for i, match := range matches {
		if len(match) < 4 {
			continue
		}
		if i >= len(indexes) || len(indexes[i]) < 8 {
			continue
		}

		prefix := line[:indexes[i][2]]
		if !hasCommentPrefix(prefix) {
			continue
		}

		text := ""
		if len(indexes[i]) >= 2 {
			text = strings.TrimSpace(line[indexes[i][1]:])
		}
		if strings.HasPrefix(text, "=") {
			continue
		}

		ref := Annotation{
			Marker: strings.ToUpper(strings.TrimSpace(match[1])),
			Path:   filepath.ToSlash(path),
			Line:   lineNumber,
			Syntax: strings.TrimSpace(match[0]),
			Text:   text,
		}

		if bracketID := strings.TrimSpace(match[2]); bracketID != "" {
			ref.TaskID = bracketID
		}
		if parenValue := strings.TrimSpace(match[3]); parenValue != "" {
			ref.Owner = strings.TrimPrefix(parenValue, "@")
		}

		refs = append(refs, ref)
	}
	return refs
}

func hasCommentPrefix(prefix string) bool {
	commentPrefixes := []string{"//", "#", "--", ";", "/*", "*", "<!--", "%", "REM ", "rem ", "@REM ", "@rem "}
	for _, commentPrefix := range commentPrefixes {
		if strings.Contains(prefix, commentPrefix) {
			return true
		}
	}
	return false
}

func isFenceLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~")
}

// resolveRepoRoot makes Patchboard work well as a submodule or standalone
// checkout inside a larger repo by finding the nearest ancestor with tasks/.
func resolveRepoRoot(start string, cfg Config) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		info, err := os.Stat(filepath.Join(current, cfg.TaskRoot))
		if err == nil && info.IsDir() {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", os.ErrNotExist
		}
		current = parent
	}
}

func parseMarkdown(body []byte) (map[string]string, string) {
	text := string(body)
	if !strings.HasPrefix(text, "---\n") {
		return map[string]string{}, text
	}

	rest := text[len("---\n"):]
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return map[string]string{}, text
	}

	return parseSimpleYAML(rest[:end]), strings.TrimPrefix(rest[end+len("\n---"):], "\n")
}

func parseSimpleYAML(text string) map[string]string {
	values := map[string]string{}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "-") || !strings.Contains(line, ":") {
			continue
		}
		key, value, _ := strings.Cut(line, ":")
		values[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
	}
	return values
}

func firstHeading(markdown string) string {
	for _, line := range strings.Split(markdown, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

func filenameSlug(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func filenameTitle(path string) string {
	slug := filenameSlug(path)
	return strings.ReplaceAll(slug, "-", " ")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

type sourcedValue struct {
	Value  string
	Source string
}

func firstNonEmptyWithSource(values []sourcedValue) (string, string) {
	for _, value := range values {
		if strings.TrimSpace(value.Value) != "" {
			return strings.TrimSpace(value.Value), value.Source
		}
	}
	return "", ""
}
