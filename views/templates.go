package views

import (
	"embed"
	"html/template"
)

//go:embed templates
var tmpl embed.FS

func MustGetTemplate() *template.Template {
	tmpl, err := template.ParseFS(tmpl, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	return tmpl
}
