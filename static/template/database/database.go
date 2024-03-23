package database

import (
	"context"

	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/migrations"
)

func Init(ctx context.Context, k *kernel.Kernel) error {
	cfg, err := di.Resolve[*config.Config](ctx, k.DependencyProvider())
	if err != nil {
		return err
	}

	return databasedi.RegisterFromConfig(
		ctx,
		k.DependencyProvider(),
		cfg.Database,
		migrations.Use(),
	)
}
