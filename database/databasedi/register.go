package databasedi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
	"github.com/jmoiron/sqlx"
)

func RegisterFromConfig[T kernel.KernelConfig](migrations *migrate.Migrations) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		di.RegisterLazySingletonWith(ctx, func(cfg T) (*sqlx.DB, error) {
			var cfgAny any = cfg
			cfger, ok := cfgAny.(dialects.DBConfiger)
			if !ok {
				return nil, fmt.Errorf("config not instance of dialects.DBConfiger")
			}
			dbcfg := cfger.DBConfig()

			dbcfg.SetDialect()
			db, err := sqlx.Open(dbcfg.DriverName(), dbcfg.DataSourceName())
			if err != nil {
				return nil, fmt.Errorf("databasedi.RegisterFromConfig: open database: %w", err)
			}

			err = migrations.Up(ctx, db)
			if err != nil {
				return nil, fmt.Errorf("databasedi.RegisterFromConfig: migrate database: %w", err)
			}
			return db, nil
		})
		registerTransactions(ctx)
		return nil
	}
}

func Register(ctx context.Context, db *sqlx.DB) {
	di.RegisterSingleton(ctx, func() *sqlx.DB {
		return db
	})
	registerTransactions(ctx)
}

func registerTransactions(ctx context.Context) {
	di.RegisterWith(ctx, func(ctx context.Context, tag string, db *sqlx.DB) (database.Read, error) {
		return func(f func(tx *sqlx.Tx) error) error {
			return runTx(ctx, db, f, true)
		}, nil
	})
	di.RegisterWith(ctx, func(ctx context.Context, tag string, db *sqlx.DB) (database.Update, error) {
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
