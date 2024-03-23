package view

import (
	"io/fs"
	"net/http"
	"text/template"
)

func Factory(fsys fs.FS, patterns ...string) func(file string) http.HandlerFunc {
	if len(patterns) == 0 {
		patterns = []string{"**/*.html"}
	}
	tpl := template.Must(template.ParseFS(fsys, patterns...))
	return func(file string) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tpl.ExecuteTemplate(w, file, nil)
		})
	}
}
