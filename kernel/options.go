package kernel

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/router"
)

type KernelOption func(*Kernel) *Kernel

func Bootstrap(bootstrap ...func(context.Context) error) KernelOption {
	return func(k *Kernel) *Kernel {
		k.bootstrap = bootstrap
		return k
	}
}

func Port(cb func() int) KernelOption {
	return func(k *Kernel) *Kernel {
		k.postBootstrap = append(k.postBootstrap, func() {
			k.port = cb()
		})
		return k
	}
}

func RootHandler(rootHandler http.Handler) KernelOption {
	return func(k *Kernel) *Kernel {
		k.rootHandler = rootHandler
		return k
	}
}

func InitRoutes(cb func(r *router.Router)) KernelOption {
	return func(k *Kernel) *Kernel {
		r := router.New()
		cb(r)
		k.rootHandler = r
		return k
	}
}

func Middleware(middleware ...router.MiddlewareFunc) KernelOption {
	return func(k *Kernel) *Kernel {
		k.middleware = middleware
		return k
	}
}

func Services(services ...Service) KernelOption {
	return func(k *Kernel) *Kernel {
		k.services = services
		return k
	}
}

func Listeners(listeners ...*Listener) KernelOption {
	return func(k *Kernel) *Kernel {
		k.listeners = map[event.EventType][]runner{}
		for _, l := range listeners {
			jobs, ok := k.listeners[l.eventType]
			if !ok {
				jobs = []runner{}
			}
			k.listeners[l.eventType] = append(jobs, l.runner)
		}
		return k
	}
}
