package dbtest

import (
	"github.com/abibby/salusa/database"
	"github.com/jmoiron/sqlx"
)

func Update(tx *sqlx.Tx) database.Update {
	return func(f func(tx *sqlx.Tx) error) error {
		return f(tx)
	}
}
func Read(tx *sqlx.Tx) database.Read {
	return func(f func(tx *sqlx.Tx) error) error {
		return f(tx)
	}
}
