package router

import (
	"fmt"
	"net/http"
	"path"

	"github.com/abibby/salusa/di"
	"github.com/gorilla/mux"
)

type MiddlewareFunc = mux.MiddlewareFunc

type WithDependencyProvider interface {
	WithDependencyProvider(dp *di.DependencyProvider)
}
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

type routeList struct {
	Routes []*Route
}

type Router struct {
	prefix string
	router *mux.Router
	dp     *di.DependencyProvider
	routes *routeList
}

func New() *Router {
	return &Router{
		prefix: "",
		router: mux.NewRouter(),
		routes: &routeList{Routes: []*Route{}},
	}
}

func (r *Router) WithDependencyProvider(dp *di.DependencyProvider) {
	r.dp = dp
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
	r.addDP(handler)
	r.router.Handle(path, handler).Methods(method)
	return r.addRoute(handler, path, method)
}
func (r *Router) Handle(p string, handler http.Handler) *Route {
	r.addDP(handler)
	r.router.PathPrefix(p).Handler(handler)
	return r.addRoute(handler, path.Join(p, "*"), "ALL")
}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.router.Use(middleware)
}

func (r *Router) Group(prefix string, cb func(r *Router)) {
	cb(&Router{
		prefix: path.Join(r.prefix, prefix),
		router: r.router.PathPrefix(prefix).Subrouter(),
		routes: r.routes,
	})
}

func (r *Router) addDP(handler http.Handler) {
	if r.dp == nil {
		return
	}
	w, ok := handler.(WithDependencyProvider)
	if !ok {
		return
	}
	w.WithDependencyProvider(r.dp)
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
