package kernel

import (
	"context"
	"net/http"
)

type KernelConfig interface {
	GetHTTPPort() int
}

type Kernel struct {
	bootstrap     []func(context.Context) error
	postBootstrap []func()
	rootHandler   func(ctx context.Context) http.Handler
	services      []Service

	cfg KernelConfig

	bootstrapped bool
}

func New[T KernelConfig](options ...KernelOption) *Kernel {
	k := &Kernel{
		bootstrap:     []func(context.Context) error{},
		postBootstrap: []func(){},
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
