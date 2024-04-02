package databasedi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

type Update func(func(tx *sqlx.Tx) error) error
type Read func(func(tx *sqlx.Tx) error) error

func RegisterFromConfig[T dialects.DBConfiger](ctx context.Context, migrations *migrate.Migrations) error {

	var dbResult *sqlx.DB
	var errResult error
	initialize := sync.OnceFunc(func() {
		cfger, err := di.Resolve[T](ctx)
		if err != nil {
			errResult = fmt.Errorf("databasedi.RegisterFromConfig: resolve config: %w", err)
			return
		}

		cfg := cfger.DBConfig()

		cfg.SetDialect()
		db, err := sqlx.Open(cfg.DriverName(), cfg.DataSourceName())
		if err != nil {
			errResult = fmt.Errorf("databasedi.RegisterFromConfig: open database: %w", err)
			return
		}

		err = migrations.Up(ctx, db)
		if err != nil {
			errResult = fmt.Errorf("databasedi.RegisterFromConfig: migrate database: %w", err)
			return
		}
		dbResult = db
	})
	di.Register(ctx, func(ctx context.Context, tag string) (*sqlx.DB, error) {
		initialize()
		return dbResult, errResult
	})
	registerTransactions(ctx)
	return nil
}

func Register(ctx context.Context, db *sqlx.DB) {
	di.RegisterSingleton(ctx, func() *sqlx.DB {
		return db
	})
	registerTransactions(ctx)
}

func registerTransactions(ctx context.Context) {
	di.RegisterWith(ctx, func(ctx context.Context, tag string, db *sqlx.DB) (Read, error) {
		return func(f func(tx *sqlx.Tx) error) error {
			return runTx(ctx, db, f, true)
		}, nil
	})
	di.RegisterWith(ctx, func(ctx context.Context, tag string, db *sqlx.DB) (Update, error) {
		return func(f func(tx *sqlx.Tx) error) error {
			return runTx(ctx, db, f, false)
		}, nil
	})
}

func runTx(ctx context.Context, db *sqlx.DB, f func(*sqlx.Tx) error, readOnly bool) error {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{
		ReadOnly: readOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
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
