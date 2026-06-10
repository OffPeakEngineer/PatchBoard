package tasks

import (
	"bytes"
	"text/template"

	pbtemplates "ledoerr/patchboard/templates"
)

func renderTemplate(name string, data any) (string, error) {
	body, err := pbtemplates.Files.ReadFile(name)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New(name).Parse(string(body))
	if err != nil {
		return "", err
	}
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", err
	}
	return rendered.String(), nil
}

func templateFile(name string) (string, error) {
	body, err := pbtemplates.Files.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
