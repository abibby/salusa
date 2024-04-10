package router

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"slices"

	"github.com/abibby/salusa/kernel"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Route struct {
	Path    string
	Method  string
	name    string
	handler http.Handler
}

func (r *Route) Name(name string) *Route {
	r.name = name
	return r
}
func (r *Route) GetName() string {
	return r.name
}
func (r *Route) GetHandler() http.Handler {
	return r.handler
}

type routeList struct {
	Routes []*Route
}

type Router struct {
	prefix     string
	middleware []MiddlewareFunc
	router     *http.ServeMux
	routes     *routeList
}

func New() *Router {
	return &Router{
		prefix:     "",
		router:     http.NewServeMux(),
		routes:     &routeList{Routes: []*Route{}},
		middleware: []MiddlewareFunc{},
	}
}

func (r *Router) Get(path string, handler http.Handler) *Route {
	return r.handleMethod(http.MethodGet, path, handler)
}
func (r *Router) Post(path string, handler http.Handler) *Route {
	return r.handleMethod(http.MethodPost, path, handler)
}
func (r *Router) Put(path string, handler http.Handler) *Route {
	return r.handleMethod(http.MethodPut, path, handler)
}
func (r *Router) Patch(path string, handler http.Handler) *Route {
	return r.handleMethod(http.MethodPatch, path, handler)
}
func (r *Router) Delete(path string, handler http.Handler) *Route {
	return r.handleMethod(http.MethodDelete, path, handler)
}

func (r *Router) handleMethod(method, pattern string, handler http.Handler) *Route {
	r.handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		handler.ServeHTTP(w, r)
	}))

	return r.addRoute(handler, pattern, method)
}
func (r *Router) Handle(pattern string, handler http.Handler) *Route {
	r.handle(pattern, handler)
	return r.addRoute(handler, pattern, "ALL")
}
func (r *Router) handle(pattern string, handler http.Handler) {
	// r.router.Handle(pattern, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// 	h := handler
	// 	for _, m := range r.middleware {
	// 		h = m(h)
	// 	}
	// 	h.ServeHTTP(w, req)
	// }))
}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.middleware = append(r.middleware, middleware)
}

func (r *Router) Group(prefix string, cb func(r *Router)) {
	cb(&Router{
		prefix:     path.Join(r.prefix, prefix),
		router:     r.router,
		routes:     r.routes,
		middleware: slices.Clone(r.middleware),
	})
}

func (r *Router) addRoute(handler http.Handler, pathName, method string) *Route {
	route := &Route{
		Path:    path.Join(r.prefix, pathName),
		Method:  method,
		handler: handler,
	}
	r.routes.Routes = append(r.routes.Routes, route)
	return route
}

func (r *Router) Routes() []*Route {
	return r.routes.Routes
}
func (r *Router) PrintRoutes() {
	for _, route := range r.routes.Routes {
		fmt.Printf("%-40s %s\n", route.Path, route.Method)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func InitRoutes(cb func(r *Router)) kernel.KernelOption {
	return kernel.RootHandler(func(ctx context.Context) http.Handler {
		r := New()
		cb(r)
		r.Register(ctx)
		return r
	})
}
