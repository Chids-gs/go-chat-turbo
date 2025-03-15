package templates

import (
	"embed"
	"html/template"
)

//go:embed views/*.html
var templateFiles embed.FS

var IndexTemplate *template.Template

func init() {
	var err error
	IndexTemplate, err = template.ParseFS(templateFiles, "views/index.html")
	if err != nil {
		panic(err)
	}
}
