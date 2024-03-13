package kernel

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/event"
)

type Kernel struct {
	bootstrap     []func(context.Context) error
	providers     []func(*di.DependencyProvider)
	postBootstrap []func()
	port          int
	rootHandler   func() http.Handler
	services      []Service
	listeners     map[event.EventType][]runner

	queue event.Queue
	dp    *di.DependencyProvider

	bootstrapped bool
}

func New(options ...KernelOption) *Kernel {
	k := &Kernel{
		bootstrap:     []func(context.Context) error{},
		providers:     []func(*di.DependencyProvider){},
		postBootstrap: []func(){},
		port:          8080,
		rootHandler: func() http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		},
		services:     []Service{},
		queue:        event.NewChannelQueue(),
		dp:           di.NewDependencyProvider(),
		bootstrapped: false,
	}

	for _, o := range options {
		k = o(k)
	}

	return k
}
