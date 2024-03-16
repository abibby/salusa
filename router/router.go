package router

import (
	"net/http"
	"path"

	"github.com/abibby/salusa/di"
	"github.com/gorilla/mux"
)

type MiddlewareFunc = mux.MiddlewareFunc

type WithDependencyProvider interface {
	WithDependencyProvider(dp *di.DependencyProvider)
}

type Router struct {
	prefix   string
	router   *mux.Router
	dp       *di.DependencyProvider
	handlers map[http.Handler]string
}

func New() *Router {
	return &Router{
		prefix:   "",
		router:   mux.NewRouter(),
		handlers: map[http.Handler]string{},
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
	r.addHandler(handler, path)
	r.router.Handle(path, handler).Methods(method)
}
func (r *Router) Handle(path string, handler http.Handler) {
	r.addDP(handler)
	r.addHandler(handler, path)
	r.router.PathPrefix(path).Handler(handler)
}

func (r *Router) Use(middleware MiddlewareFunc) {
	r.router.Use(middleware)
}

func (r *Router) Group(prefix string, cb func(r *Router)) {
	cb(&Router{
		prefix: path.Join(r.prefix, prefix),
		router: r.router.PathPrefix(prefix).Subrouter(),
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
func (r *Router) addHandler(handler http.Handler, path string) {
	r.handlers[handler] = path
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
