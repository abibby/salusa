package router

import (
	"net/http"

	"github.com/abibby/salusa/openapidoc"
	"github.com/go-openapi/spec"
)

var _ openapidoc.Pathser = (*Router)(nil)

func (r *Router) Paths() (*spec.Paths, error) {
	var err error
	paths := &spec.Paths{
		Paths: map[string]spec.PathItem{},
	}
	for _, route := range r.routes.Routes {
		pathItem, ok := paths.Paths[route.Path]
		if !ok {
			pathItem = spec.PathItem{}
		}

		var op *spec.Operation
		if oper, ok := route.handler.(openapidoc.Operationer); ok {
			op, err = oper.Operation()
			if err != nil {
				return nil, err
			}
		} else {
			// op = spec.NewOperation("")
		}

		switch route.Method {
		case http.MethodGet:
			pathItem.Get = op
		case http.MethodHead:
			pathItem.Head = op
		case http.MethodPost:
			pathItem.Post = op
		case http.MethodPut:
			pathItem.Put = op
		case http.MethodPatch:
			pathItem.Patch = op
		case http.MethodDelete:
			pathItem.Delete = op
		case http.MethodOptions:
			pathItem.Options = op
		default:
			continue
			// return nil, fmt.Errorf("unsupported method %s", route.Method)
		}
		paths.Paths[route.Path] = pathItem
	}
	return paths, nil
}
