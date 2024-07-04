package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
}

type Update func(func(tx *sqlx.Tx) error) error
type Read func(func(tx *sqlx.Tx) error) error

func Exec(ctx context.Context, db DB, query string, args []interface{}) (sql.Result, error) {
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("SQL Exec: %s: %w", query, err)
	}
	return result, nil
}

func runTx(ctx context.Context, db DB, f func(*sqlx.Tx) error, readOnly bool) error {
	var tx *sqlx.Tx
	var err error
	switch db := db.(type) {
	case *sqlx.DB:
		tx, err = db.BeginTxx(ctx, &sql.TxOptions{
			ReadOnly: readOnly,
		})
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}
	case *sqlx.Tx:
		tx = db
	default:
		return fmt.Errorf("unsupported type %v", reflect.TypeOf(db))
	}

	defer func() {
		err := recover()
		if err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()

	err = f(tx)
	if err == nil {
		txErr := tx.Commit()
		if txErr != nil {
			return txErr
		}
	} else {
		txErr := tx.Rollback()
		if txErr != nil {
			return errors.Join(err, txErr)
		}
	}
	return err
}
func NewUpdate(ctx context.Context, db DB) Update {
	return func(f func(tx *sqlx.Tx) error) (err error) {
		return runTx(ctx, db, f, false)
	}
}
func NewRead(ctx context.Context, db DB) Read {
	return func(f func(tx *sqlx.Tx) error) error {
		return runTx(ctx, db, f, true)
	}
}
