package kernel

import (
	"context"
	"net/http"

	"github.com/go-openapi/spec"
)

type KernelConfig interface {
	GetHTTPPort() int
	GetBaseURL() string
}

type Kernel struct {
	bootstrap      []func(context.Context) error
	registerConfig func(context.Context) error
	rootHandler    func(ctx context.Context) http.Handler
	services       []Service

	docs *spec.Swagger

	cfg KernelConfig

	bootstrapped bool
}

func New(options ...KernelOption) *Kernel {
	k := &Kernel{
		bootstrap:      []func(context.Context) error{},
		registerConfig: func(context.Context) error { return nil },
		rootHandler: func(ctx context.Context) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		},
		services:     []Service{},
		bootstrapped: false,
	}

	for _, o := range options {
		k = o(k)
	}

	return k
}

func (k *Kernel) RootHandler(ctx context.Context) http.Handler {
	return k.rootHandler(ctx)
}

func (k *Kernel) Config() KernelConfig {
	return k.cfg
}
