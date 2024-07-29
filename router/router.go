package router

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

type Middleware interface {
	Middleware(next http.Handler) http.Handler
}
type MiddlewareFunc func(http.Handler) http.Handler

func (f MiddlewareFunc) Middleware(next http.Handler) http.Handler {
	return f(next)
}

type Route struct {
	Path    string
	Method  string
	name    string
	handler http.Handler
	router  *Router
}

func (r *Route) GetMiddleware() []Middleware {
	return r.router.middleware
}

func (r *Route) Name(name string) *Route {
	r.name = name
	return r
}

type routeList struct {
	Routes []*Route
}

type Router struct {
	prefix     string
	router     *mux.Router
	routes     *routeList
	middleware []Middleware
}

func New() *Router {
	return &Router{
		prefix:     "",
		router:     mux.NewRouter(),
		routes:     &routeList{Routes: []*Route{}},
		middleware: []Middleware{},
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

func (r *Router) handleMethod(method, path string, handler http.Handler) *Route {
	r.router.Handle(path, handler).Methods(method)
	return r.addRoute(handler, path, method)
}
func (r *Router) Handle(p string, handler http.Handler) *Route {
	r.router.PathPrefix(p).Handler(handler)
	return r.addRoute(handler, path.Join(p, "*"), "ALL")
}

func (r *Router) UseFunc(middleware func(http.Handler) http.Handler) {
	m := MiddlewareFunc(middleware)
	r.router.Use(m.Middleware)
	r.middleware = append(r.middleware, m)
}
func (r *Router) Use(middleware Middleware) {
	r.router.Use(middleware.Middleware)
	r.middleware = append(r.middleware, middleware)
}

func (r *Router) Group(prefix string, cb func(r *Router)) {
	middleware := make([]Middleware, len(r.middleware))
	copy(middleware, r.middleware)
	cb(&Router{
		prefix:     path.Join(r.prefix, prefix),
		router:     r.router.PathPrefix(prefix).Subrouter(),
		routes:     r.routes,
		middleware: middleware,
	})
}

func (r *Router) addRoute(handler http.Handler, pathName, method string) *Route {
	route := &Route{
		Path:    path.Join(r.prefix, pathName),
		Method:  method,
		handler: handler,
		router:  r,
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
