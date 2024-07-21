package openapidocdi

import (
	"context"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/openapidoc"
)

func Register(ctx context.Context) error {
	di.RegisterLazySingletonWith(ctx, func(k *kernel.Kernel) (openapidoc.APIDocer, error) {
		return k, nil
	})
	return nil
}
