package databasedi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

type Update func(func(tx *sqlx.Tx) error) error
type Read func(func(tx *sqlx.Tx) error) error

func RegisterFromConfig(ctx context.Context, cfg dialects.Config, migrations *migrate.Migrations) error {
	cfg.SetDialect()
	db, err := sqlx.Open(cfg.DriverName(), cfg.DataSourceName())
	if err != nil {
		return err
	}

	err = migrations.Up(ctx, db)
	if err != nil {
		return err
	}

	Register(ctx, db)
	return nil
}

func Register(ctx context.Context, db *sqlx.DB) {
	di.RegisterSingleton(ctx, func() *sqlx.DB {
		return db
	})
	di.Register(ctx, func(ctx context.Context, tag string) (Read, error) {
		return func(f func(tx *sqlx.Tx) error) error {
			return runTx(ctx, db, f, true)
		}, nil
	})
	di.Register(ctx, func(ctx context.Context, tag string) (Update, error) {
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
