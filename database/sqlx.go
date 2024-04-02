package database

import "github.com/jmoiron/sqlx"

type DB interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
}

type Update func(func(tx *sqlx.Tx) error) error
type Read func(func(tx *sqlx.Tx) error) error
