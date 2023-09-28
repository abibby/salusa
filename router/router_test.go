package router_test

import (
	"testing"

	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

func ExampleRouter(t *testing.T) {
	r := router.New()

	r.Group("/test", func(r *router.Router) {
		r.Get("/", request.Handler(func(r *any) (any, error) {
			return nil, nil
		}))
	})
}
