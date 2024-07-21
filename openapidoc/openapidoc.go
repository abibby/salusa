package openapidoc

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/abibby/salusa/di"
	"github.com/go-openapi/spec"
)

type Operationer interface {
	Operation(ctx context.Context) (*spec.Operation, error)
}
type Pathser interface {
	Paths(ctx context.Context, basePath string) (*spec.Paths, error)
}
type APIDocer interface {
	APIDoc(context.Context) (*spec.Swagger, error)
}

//go:generate npm run upgrade-and-install

//go:embed node_modules/swagger-ui-dist
var embededFiles embed.FS

//go:embed swagger-initializer.js
var swaggerInitializer []byte

func SwaggerUI(prefix string) http.Handler {
	return http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		swaggerUIFiles, err := fs.Sub(embededFiles, "node_modules/swagger-ui-dist")
		if err != nil {
			panic(err)
		}
		switch r.URL.Path {
		case "":
			http.Redirect(w, r, r.RequestURI+"/", http.StatusFound)

		case "/swagger-initializer.js":
			w.Header().Add("Content-Length", fmt.Sprint(len(swaggerInitializer)))
			w.Header().Add("Content-Type", "application/json")
			w.Write(swaggerInitializer)

		case "/swagger.json":
			api, err := di.Resolve[APIDocer](r.Context())
			if err != nil {
				panic(err)
			}
			doc, err := api.APIDoc(r.Context())
			if err != nil {
				panic(err)
			}
			e := json.NewEncoder(w)
			e.SetIndent("", "    ")
			err = e.Encode(doc)
			if err != nil {
				panic(err)
			}

		default:
			http.FileServerFS(swaggerUIFiles).ServeHTTP(w, r)
		}
	}))
}
