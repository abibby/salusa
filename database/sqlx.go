package database

import "github.com/jmoiron/sqlx"

type DB interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
}
