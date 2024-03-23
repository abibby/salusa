package router_test

import (
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

func ExampleRouter() {
	r := router.New()

	r.Group("/test", func(r *router.Router) {
		r.Get("/", request.Handler(func(r *any) (any, error) {
			return nil, nil
		}))
	})
}
