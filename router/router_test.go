package router_test

import (
	"abibby.com/salusa/request"
	"abibby.com/salusa/router"
)

func ExampleRouter() {
	r := router.New()

	r.Group("/test", func(r *router.Router) {
		r.Get("/", request.Handler(func(r *any) (any, error) {
			return nil, nil
		}))
	})
}
