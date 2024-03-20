package router

import (
	"bytes"
	"fmt"
	"net/http"
	"path"
	"reflect"

	"github.com/abibby/salusa/di"
	"github.com/gorilla/mux"
)

type MiddlewareFunc = mux.MiddlewareFunc

type WithDependencyProvider interface {
	WithDependencyProvider(dp *di.DependencyProvider)
}
type Route struct {
	Path   string
	Method string
}
type routeList struct {
	Routes []*Route
}

type Router struct {
	prefix   string
	router   *mux.Router
	dp       *di.DependencyProvider
	handlers map[http.Handler]string
	routes   *routeList
}

func New() *Router {
	return &Router{
		prefix:   "",
		router:   mux.NewRouter(),
		handlers: map[http.Handler]string{},
		routes:   &routeList{Routes: []*Route{}},
	}
}

func (r *Router) WithDependencyProvider(dp *di.DependencyProvider) {
	r.dp = dp
}

func (r *Router) Get(path string, handler http.Handler) {
	r.handleMethod(http.MethodGet, path, handler)
}
func (r *Router) Post(path string, handler http.Handler) {
	r.handleMethod(http.MethodPost, path, handler)
}
func (r *Router) Put(path string, handler http.Handler) {
	r.handleMethod(http.MethodPut, path, handler)
}
func (r *Router) Patch(path string, handler http.Handler) {
	r.handleMethod(http.MethodPatch, path, handler)
}
func (r *Router) Delete(path string, handler http.Handler) {
	r.handleMethod(http.MethodDelete, path, handler)
}

func (r *Router) handleMethod(method, path string, handler http.Handler) {
	r.addDP(handler)
	r.addHandler(handler, path, method)
	r.router.Handle(path, handler).Methods(method)
}
func (r *Router) Handle(p string, handler http.Handler) {
	r.addDP(handler)
	r.addHandler(handler, path.Join(p, "*"), "ALL")
	r.router.PathPrefix(p).Handler(handler)
}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.router.Use(middleware)
}

func (r *Router) Group(prefix string, cb func(r *Router)) {
	cb(&Router{
		prefix:   path.Join(r.prefix, prefix),
		router:   r.router.PathPrefix(prefix).Subrouter(),
		handlers: r.handlers,
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
func (r *Router) addHandler(handler http.Handler, p, m string) {
	if reflect.ValueOf(handler).Comparable() {
		r.handlers[handler] = path.Join(r.prefix, p)
	}
	r.routes.Routes = append(r.routes.Routes, &Route{
		Path:   p,
		Method: m,
	})
}

func (r *Router) Routes() string {
	b := &bytes.Buffer{}
	for _, route := range r.routes.Routes {
		fmt.Fprintf(b, "%-40s %s\n", route.Path, route.Method)
	}
	return b.String()
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
