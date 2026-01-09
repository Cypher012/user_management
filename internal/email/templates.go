package email

import (
	"bytes"
	"html/template"
)

type Template struct {
	verify  *template.Template
	reset   *template.Template
	welcome *template.Template
}

func LoadTemplates() (*Template, error) {
	verify, err := template.ParseFiles("internal/email/templates/verify.html")
	if err != nil {
		return nil, err
	}
	reset, err := template.ParseFiles("internal/email/templates/reset.html")
	if err != nil {
		return nil, err
	}

	welcome, err := template.ParseFiles("internal/email/templates/welcome.html")
	if err != nil {
		return nil, err
	}

	return &Template{
		verify:  verify,
		reset:   reset,
		welcome: welcome,
	}, nil
}

func render(t *template.Template, data any) (string, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
