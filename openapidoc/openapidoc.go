package openapidoc

import (
	"embed"
	"io/fs"
	"net/http"
	"text/template"

	"github.com/go-openapi/spec"
)

type Operationer interface {
	Operation() (*spec.Operation, error)
}
type Pathser interface {
	Paths() (*spec.Paths, error)
}
type APIDocer interface {
	APIDoc() (*spec.Swagger, error)
}

// not go:generate npm upgrade && npm i

//go:embed node_modules/swagger-ui-dist
var embededFiles embed.FS

//go:embed swagger-initializer.js
var swaggerInit string

func SwaggerUI(url string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		swaggerUIFiles, err := fs.Sub(embededFiles, "node_modules/swagger-ui-dist")
		if err != nil {
			panic(err)
		}
		if r.URL.Path == "/swagger-initializer.js" {
			tpl, err := template.New("").Parse(swaggerInit)
			if err != nil {
				panic(err)
			}
			tpl.Execute(w, map[string]string{"URL": url})
			return
		}
		http.FileServerFS(swaggerUIFiles).ServeHTTP(w, r)
	})
}
