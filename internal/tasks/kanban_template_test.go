package tasks

import (
	"regexp"
	"strings"
	"testing"
)

func TestKanbanTranslationsHaveMatchingUniqueKeys(t *testing.T) {
	body, err := templateFile("kanban.html")
	if err != nil {
		t.Fatalf("loading kanban template: %v", err)
	}

	blockPattern := regexp.MustCompile(`(?s)(en|es): \{(.*?)\n      \}(?:,|\n)`)
	keyPattern := regexp.MustCompile(`(?m)^        ([A-Za-z][A-Za-z0-9]*):`)
	matches := blockPattern.FindAllStringSubmatch(body, -1)
	if len(matches) != 2 {
		t.Fatalf("expected English and Spanish message blocks, found %d", len(matches))
	}

	keysByLanguage := map[string]map[string]bool{}
	for _, match := range matches {
		keys := map[string]bool{}
		for _, keyMatch := range keyPattern.FindAllStringSubmatch(match[2], -1) {
			key := keyMatch[1]
			if keys[key] {
				t.Fatalf("duplicate %s translation key %q", match[1], key)
			}
			keys[key] = true
		}
		keysByLanguage[match[1]] = keys
	}

	for key := range keysByLanguage["en"] {
		if !keysByLanguage["es"][key] {
			t.Errorf("Spanish messages are missing key %q", key)
		}
	}
	for key := range keysByLanguage["es"] {
		if !keysByLanguage["en"][key] {
			t.Errorf("English messages are missing key %q", key)
		}
	}
}

func TestKanbanTemplateIncludesRememberedFolderStateMachine(t *testing.T) {
	body, err := templateFile("kanban.html")
	if err != nil {
		t.Fatalf("loading kanban template: %v", err)
	}

	for _, required := range []string{
		`indexedDB.open("patchboard", 1)`,
		`queryPermission({ mode })`,
		`requestPermission({ mode })`,
		`id: "patchboard-tasks", mode: "readwrite"`,
		`reconnectFolderAction`,
		`changeFolderAction`,
		`forgetFolder`,
		`error.name === "AbortError"`,
		`error.code = "PATCHBOARD_WRONG_FOLDER"`,
		`installEl.hidden = !installPrompt`,
		`el.draggable = accessMode === "writable"`,
	} {
		if !strings.Contains(body, required) {
			t.Errorf("kanban template is missing %q", required)
		}
	}
}
