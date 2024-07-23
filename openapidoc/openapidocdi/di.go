package openapidocdi

import (
	"context"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/openapidoc"
)

type apiDocerOpts struct {
	Kernel *kernel.Kernel `inject:""`
}

func Register(ctx context.Context) error {
	di.RegisterLazySingletonWith(ctx, func(opts *apiDocerOpts) (openapidoc.APIDocer, error) {
		return opts.Kernel, nil
	})
	return nil
}
