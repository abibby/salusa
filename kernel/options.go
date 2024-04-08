package kernel

import (
	"context"
	"net/http"

	"github.com/abibby/salusa/di"
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
			di.RegisterSingleton(ctx, func() KernelConfig {
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

func Services(services ...Service) KernelOption {
	return func(k *Kernel) *Kernel {
		k.services = services
		return k
	}
}
