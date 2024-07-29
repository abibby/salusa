package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/jmoiron/sqlx"
)

type DB interface {
	sqlx.QueryerContext
	sqlx.ExecerContext
}

type Transaction interface {
	Run(func(tx *sqlx.Tx) error) error
}

type Update func(func(tx *sqlx.Tx) error) error

func (f Update) Run(cb func(tx *sqlx.Tx) error) error {
	return f(cb)
}

type Read func(func(tx *sqlx.Tx) error) error

func (f Read) Run(cb func(tx *sqlx.Tx) error) error {
	return f(cb)
}

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
func NewUpdate(ctx context.Context, mtx sync.Locker, db DB) Update {
	return func(f func(tx *sqlx.Tx) error) (err error) {
		if mtx != nil {
			mtx.Lock()
			defer mtx.Unlock()
		}
		return runTx(ctx, db, f, false)
	}
}
func NewRead(ctx context.Context, mtx sync.Locker, db DB) Read {
	return func(f func(tx *sqlx.Tx) error) error {
		if mtx != nil {
			if rwmtx, ok := mtx.(*sync.RWMutex); ok {
				rwmtx.RLock()
				defer rwmtx.RUnlock()
			} else {
				mtx.Lock()
				defer mtx.Unlock()
			}
		}
		return runTx(ctx, db, f, true)
	}
}

func Value[T any](txFunc Transaction, cb func(tx *sqlx.Tx) (T, error)) (T, error) {
	var result T
	err := txFunc.Run(func(tx *sqlx.Tx) error {
		var err error
		result, err = cb(tx)
		return err
	})
	return result, err
}
