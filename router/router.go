package router

import (
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

type MiddlewareFunc = mux.MiddlewareFunc

type Router struct {
	prefix string
	router *mux.Router
}

func New() *Router {
	return &Router{
		prefix: "",
		router: mux.NewRouter(),
	}
}

func (r *Router) Get(path string, handler http.Handler) {
	r.router.Handle(path, handler).Methods(http.MethodGet)
}
func (r *Router) Post(path string, handler http.Handler) {
	r.router.Handle(path, handler).Methods(http.MethodPost)
}
func (r *Router) Put(path string, handler http.Handler) {
	r.router.Handle(path, handler).Methods(http.MethodPut)
}
func (r *Router) Patch(path string, handler http.Handler) {
	r.router.Handle(path, handler).Methods(http.MethodPatch)
}
func (r *Router) Delete(path string, handler http.Handler) {
	r.router.Handle(path, handler).Methods(http.MethodDelete)
}

func (r *Router) Handle(path string, handler http.Handler) {
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

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}
