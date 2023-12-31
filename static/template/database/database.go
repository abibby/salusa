package database

import (
	"context"

	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/migrations"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

func Init(ctx context.Context) error {
	sqlite.UseSQLite()
	db, err := sqlx.Open("sqlite", config.DBPath)
	if err != nil {
		return err
	}

	err = migrations.Use().Up(ctx, db)
	if err != nil {
		return err
	}

	DB = db

	return nil
}

func WithDB() router.MiddlewareFunc {
	return request.WithDB(DB)
}
