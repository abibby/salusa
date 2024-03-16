package kernel

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/router"
)

type KernelOption func(*Kernel) *Kernel

type KernelConfig struct {
	Port int
}

func Bootstrap(bootstrap ...func(context.Context, *Kernel) error) KernelOption {
	return func(k *Kernel) *Kernel {
		k.bootstrap = bootstrap
		return k
	}
}

func Config(cb func() *KernelConfig) KernelOption {
	return func(k *Kernel) *Kernel {
		k.postBootstrap = append(k.postBootstrap, func() {
			cfg := cb()
			k.port = cfg.Port
		})
		return k
	}
}

func RootHandler(rootHandler func() http.Handler) KernelOption {
	return func(k *Kernel) *Kernel {
		k.rootHandler = rootHandler
		return k
	}
}

func InitRoutes(cb func(r *router.Router)) KernelOption {
	return func(k *Kernel) *Kernel {
		k.rootHandler = func() http.Handler {
			r := router.New()
			r.WithDependencyProvider(k.dp)
			cb(r)
			k.Register(r.Register)
			return r
		}
		return k
	}
}

func Services(services ...Service) KernelOption {
	return func(k *Kernel) *Kernel {
		k.services = services
		return k
	}
}

// func Listeners(listeners ...*Listener) KernelOption {
// 	return func(k *Kernel) *Kernel {
// 		k.listeners = map[event.EventType][]runner{}
// 		for _, l := range listeners {
// 			jobs, ok := k.listeners[l.eventType]
// 			if !ok {
// 				jobs = []runner{}
// 			}
// 			k.listeners[l.eventType] = append(jobs, l.runner)
// 		}
// 		return k
// 	}
// }

func Providers(providers ...func(*di.DependencyProvider)) KernelOption {
	return func(k *Kernel) *Kernel {
		k.providers = providers
		return k
	}
}
