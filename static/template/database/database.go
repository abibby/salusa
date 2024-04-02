package database

import (
	"context"

	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/migrations"
)

func Init(ctx context.Context) error {
	return databasedi.RegisterFromConfig[*config.Config](
		ctx,
		migrations.Use(),
	)
}
