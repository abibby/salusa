package kernel

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/router"
)

type Kernel struct {
	bootstrap     []func(context.Context) error
	postBootstrap []func()
	port          int
	rootHandler   http.Handler
	middleware    []router.MiddlewareFunc
	services      []Service
	listeners     map[event.EventType][]runner
	queue         event.Queue
}

var defaultKernel *Kernel

func New(options ...KernelOption) *Kernel {
	k := &Kernel{
		bootstrap:     []func(context.Context) error{},
		postBootstrap: []func(){},
		port:          8080,
		rootHandler:   http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
		middleware:    []router.MiddlewareFunc{},
		services:      []Service{},
		queue:         event.NewChannelQueue(),
	}

	for _, o := range options {
		k = o(k)
	}

	return k
}

func NewDefaultKernel(options ...KernelOption) *Kernel {
	defaultKernel = New(options...)
	return defaultKernel
}
