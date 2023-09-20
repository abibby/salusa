package kernel

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/event"
)

type Kernel struct {
	bootstrap     []func(context.Context) error
	postBootstrap []func()
	port          int
	rootHandler   func() http.Handler
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
		rootHandler: func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		},
		services: []Service{},
		queue:    event.NewChannelQueue(),
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
