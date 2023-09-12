package kernel

import (
	"net/http"

	"github.com/abibby/salusa/router"
)

type KernelOption func(*Kernel) *Kernel

func Bootstrap(bootstrap ...func() error) KernelOption {
	return func(k *Kernel) *Kernel {
		k.bootstrap = bootstrap
		return k
	}
}

func Port(port int) KernelOption {
	return func(k *Kernel) *Kernel {
		k.port = port
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
