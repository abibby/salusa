package view

import (
	"bytes"
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

type viewData struct {
	URL router.URLResolver `inject:""`
}

type ViewFunc func(file string) http.HandlerFunc

func Factory(dp *di.DependencyProvider, fsys fs.FS, patterns ...string) ViewFunc {
	if len(patterns) == 0 {
		patterns = []string{"**/*.html"}
	}

	tpl, err := template.ParseFS(fsys, patterns...)
	if err != nil {
		panic(err)
	}

	return func(file string) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := &viewData{}
			err := dp.Fill(r.Context(), data)
			if err != nil {
				request.NewHTTPError(err, http.StatusInternalServerError).Respond(w, r)
				return
			}
			b := &bytes.Buffer{}
			err = tpl.Funcs(template.FuncMap{
				"route": data.URL.Resolve,
			}).ExecuteTemplate(b, file, data)
			if err != nil {
				request.NewHTTPError(err, http.StatusInternalServerError).Respond(w, r)
				return
			}
			io.Copy(w, b)
		})
	}
}
