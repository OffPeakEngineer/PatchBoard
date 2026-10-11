// Copyright (c) 2026 Andrew David LeTourneau; MIT OR Zlib
package tasks

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestSummaryExportsPortableBoardInAllFormats(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", "states: [ready, done]\ndone_states: [done]\n")
	if err := os.Mkdir(filepath.Join(root, "tasks", "done"), 0o755); err != nil {
		t.Fatal(err)
	}
	markdown := "---\nid: task-launch\n---\n# Launch, \"beta\" <script>alert(1)</script>\n\nFull details <img src=x onerror=alert(1)> & more.\n"
	writeFile(t, root, "tasks/ready/launch.md", markdown)
	writeFile(t, root, "tasks/unconfigured/follow-up.md", "# Follow up\n")
	writeFile(t, root, "tasks/templates/example.md", "# Template\n")

	for _, format := range []string{"yaml", "json", "csv", "html"} {
		t.Run(format, func(t *testing.T) {
			result, err := ExportSummary(filepath.Join(root, "tasks"), SummaryOptions{Type: format})
			if err != nil {
				t.Fatal(err)
			}
			if result.Tasks != 2 || result.Path != filepath.ToSlash(filepath.Join(root, "summary."+format)) {
				t.Fatalf("unexpected export result: %#v", result)
			}
			body, err := os.ReadFile(result.Path)
			if err != nil {
				t.Fatal(err)
			}
			switch format {
			case "yaml", "json":
				var board BoardSummary
				if format == "yaml" {
					err = yaml.Unmarshal(body, &board)
				} else {
					err = json.Unmarshal(body, &board)
				}
				if err != nil {
					t.Fatal(err)
				}
				if board.TaskRoot != "tasks" || board.TotalTasks != 2 || len(board.States) != 3 {
					t.Fatalf("incorrect board: %#v", board)
				}
				if board.States[0].State != "ready" || board.States[0].Tasks[0].ID != "task-launch" ||
					board.States[1].State != "done" || board.States[1].Count != 0 || board.States[1].Tasks == nil ||
					board.States[2].State != "unconfigured" {
					t.Fatalf("incorrect state ordering or contents: %#v", board.States)
				}
				if bytes.Contains(body, []byte("Full details")) {
					t.Fatal("structured summaries should contain task metadata, not Markdown bodies")
				}
			case "csv":
				rows, err := csv.NewReader(bytes.NewReader(body)).ReadAll()
				if err != nil {
					t.Fatal(err)
				}
				if len(rows) != 3 || strings.Join(rows[0], ",") != "id,title,state,path" ||
					rows[1][1] != "Launch, \"beta\" <script>alert(1)</script>" || rows[1][3] != "tasks/ready/launch.md" {
					t.Fatalf("incorrect CSV escaping or rows: %#v", rows)
				}
			case "html":
				page := string(body)
				if strings.Count(page, `class="lane"`) != 3 || strings.Count(page, `class="card"`) != 2 ||
					!strings.Contains(page, `href="#task-`) || !strings.Contains(page, `id="task-`) ||
					!strings.Contains(page, "Full details &lt;img") || !strings.Contains(page, "&lt;script&gt;") {
					t.Fatalf("missing offline board contents: %s", page)
				}
				contents := regexp.MustCompile(`(?s)<pre>(.*?)</pre>`).FindAllStringSubmatch(page, -1)
				if len(contents) != 2 || html.UnescapeString(contents[0][1]) != markdown {
					t.Fatal("embedded task text does not preserve the original Markdown")
				}
				anchors := map[string]bool{}
				for _, match := range regexp.MustCompile(`\bid="([^"]+)"`).FindAllStringSubmatch(page, -1) {
					if anchors[match[1]] {
						t.Fatalf("duplicate HTML anchor %q", match[1])
					}
					anchors[match[1]] = true
				}
				for _, match := range regexp.MustCompile(`href="#([^"]+)"`).FindAllStringSubmatch(page, -1) {
					if !anchors[match[1]] {
						t.Fatalf("offline link has no destination: %q", match[1])
					}
				}
				for _, unsafe := range []string{"<script", "<img", "<button", "<select", "https://", "showDirectoryPicker"} {
					if strings.Contains(page, unsafe) {
						t.Fatalf("HTML contains active or external content %q", unsafe)
					}
				}
			}
			if _, err := ExportSummary(root, SummaryOptions{Type: format}); err != nil {
				t.Fatal(err)
			}
			again, err := os.ReadFile(result.Path)
			if err != nil || !bytes.Equal(body, again) {
				t.Fatalf("repeated export changed snapshot: %v", err)
			}
		})
	}
}

func TestSummaryDefaultsToYAMLAndAllowsCustomOutput(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/ready/task.md", "# Task\n")
	writeFile(t, root, "tasks/board.yml", "states: [ready]\n")
	output := filepath.Join(root, "reports", "board.yaml")
	result, err := ExportSummary(root, SummaryOptions{Output: output})
	if err != nil || result.Path != filepath.ToSlash(output) {
		t.Fatalf("custom export = %#v, %v", result, err)
	}
	body, err := os.ReadFile(output)
	if err != nil || !bytes.Contains(body, []byte("total_tasks: 1")) {
		t.Fatalf("incorrect default YAML: %s, %v", body, err)
	}
}

func TestSummaryRejectsInvalidFormatsAndProtectsBoardFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/ready/task.md", "# Task\n")
	writeFile(t, root, "tasks/board.yml", "states: [ready]\n")
	writeFile(t, root, "tasks/kanban.html", "editable board")
	for _, opts := range []SummaryOptions{
		{Type: "pdf"},
		{Type: "yaml", Output: filepath.Join(root, "tasks", "board.yml")},
		{Type: "html", Output: filepath.Join(root, "tasks", "kanban.html")},
		{Type: "json", Output: filepath.Join(root, "tasks", "ready", "task.md")},
	} {
		if _, err := ExportSummary(root, opts); err == nil {
			t.Fatalf("expected invalid export to fail: %#v", opts)
		}
	}
	body, err := os.ReadFile(filepath.Join(root, "tasks", "ready", "task.md"))
	if err != nil || string(body) != "# Task\n" {
		t.Fatalf("export changed task: %s, %v", body, err)
	}
}

func TestSummaryExportsEmptyBoardWithoutKanban(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/board.yml", "states: [ready]\n")
	if err := os.Mkdir(filepath.Join(root, "tasks", "ready"), 0o755); err != nil {
		t.Fatal(err)
	}
	result, err := ExportSummary(root, SummaryOptions{Type: "html"})
	if err != nil || result.Tasks != 0 {
		t.Fatalf("empty board export = %#v, %v", result, err)
	}
	writeFile(t, root, "tasks/board.yml", "states: [missing]\n")
	if _, err := ExportSummary(root, SummaryOptions{Type: "html"}); err == nil {
		t.Fatal("expected unrelated folders to fail export")
	}
}

func TestHTMLSummaryDoesNotIntroduceCopiedCodeAnnotations(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "tasks/2_doing/example.md", "---\nid: task-example\n---\n# Example\n\n// "+todoRef("task-example")+": example annotation\n")
	if _, err := ExportSummary(root, SummaryOptions{Type: "html"}); err != nil {
		t.Fatal(err)
	}
	result, err := Lint(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Todos) != 0 {
		t.Fatalf("exported task text introduced code annotations: %#v", result.Todos)
	}
}
