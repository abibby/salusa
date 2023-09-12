package database

import (
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/config"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var DB *sqlx.DB

func Init() error {
	db, err := sqlx.Open("sqlite", config.DBPath)
	if err != nil {
		return err
	}
	DB = db

	return nil
}

func WithDB() router.MiddlewareFunc {
	return request.WithDB(DB)
}
