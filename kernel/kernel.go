package kernel

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/config"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/go-openapi/spec"
)

type Kernel struct {
	bootstrap          []func(context.Context) error
	registerConfig     func(context.Context) error
	rootHandlerFactory func(ctx context.Context) http.Handler
	rootHandler        http.Handler
	services           []Service

	globalMiddleware []router.MiddlewareFunc

	docs *spec.Swagger

	cfg config.Config

	bootstrapped bool
}

func New(options ...KernelOption) *Kernel {
	k := &Kernel{
		bootstrap:      []func(context.Context) error{},
		registerConfig: func(context.Context) error { return nil },
		rootHandler:    http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		globalMiddleware: []router.MiddlewareFunc{
			request.DIMiddleware(),
		},
		services:     []Service{},
		bootstrapped: false,
	}

	for _, o := range options {
		k = o(k)
	}

	return k
}

func (k *Kernel) RootHandler() http.Handler {
	if !k.bootstrapped {
		panic("cannot access root handler before the kernel is bootstrapped")
	}
	return k.rootHandler
}

func (k *Kernel) Config() config.Config {
	return k.cfg
}
