package database

import (
	"context"

	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/migrations"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Init(ctx context.Context, k *kernel.Kernel) error {
	sqlite.UseSQLite()
	db, err := sqlx.Open("sqlite", config.DBPath)
	if err != nil {
		return err
	}

	err = migrations.Use().Up(ctx, db)
	if err != nil {
		return err
	}

	databasedi.Register(k.DependencyProvider(), db)

	return nil
}
