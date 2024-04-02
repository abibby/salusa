package kernel

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/router"
)

type KernelOption func(*Kernel) *Kernel

func Bootstrap(bootstrap ...func(context.Context) error) KernelOption {
	return func(k *Kernel) *Kernel {
		k.bootstrap = bootstrap
		return k
	}
}

func Config[T KernelConfig](cb func() T) KernelOption {
	return func(k *Kernel) *Kernel {
		cfg := cb()
		k.cfg = cfg
		k.postBootstrap = append(k.postBootstrap, func(ctx context.Context) error {
			di.RegisterSingleton(ctx, func() T {
				return cfg
			})
			return nil
		})
		return k
	}
}

func RootHandler(rootHandler func(ctx context.Context) http.Handler) KernelOption {
	return func(k *Kernel) *Kernel {
		k.rootHandler = rootHandler
		return k
	}
}

func InitRoutes(cb func(r *router.Router)) KernelOption {
	return func(k *Kernel) *Kernel {
		k.rootHandler = func(ctx context.Context) http.Handler {
			r := router.New()
			cb(r)
			r.Register(ctx)
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
