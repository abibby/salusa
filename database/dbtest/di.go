package dbtest

import (
	"github.com/abibby/salusa/database/databasedi"
	"github.com/jmoiron/sqlx"
)

func Update(tx *sqlx.Tx) databasedi.Update {
	return func(f func(tx *sqlx.Tx) error) error {
		return f(tx)
	}
}
func Read(tx *sqlx.Tx) databasedi.Read {
	return func(f func(tx *sqlx.Tx) error) error {
		return f(tx)
	}
}
