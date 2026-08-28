package openapidocdi

import (
	"context"

	"abibby.com/salusa/di"
	"abibby.com/salusa/kernel"
	"abibby.com/salusa/openapidoc"
)

type apiDocerOpts struct {
	Kernel *kernel.Kernel `inject:""`
}

func Register(ctx context.Context) {
	di.RegisterLazySingletonWith(ctx, func(opts *apiDocerOpts) (openapidoc.APIDocer, error) {
		return opts.Kernel, nil
	})
}
