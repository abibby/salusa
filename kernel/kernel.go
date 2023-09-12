package kernel

import (
	"net/http"

	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/config"
)

type Kernel struct {
	bootstrap   []func() error
	port        int
	rootHandler http.Handler
	middleware  []router.MiddlewareFunc
}

func New(options ...KernelOption) *Kernel {
	k := &Kernel{
		bootstrap:   []func() error{},
		port:        config.Port,
		rootHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		middleware:  []router.MiddlewareFunc{},
	}

	for _, o := range options {
		k = o(k)
	}

	return k
}
