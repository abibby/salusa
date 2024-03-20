package databasedi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

func Register(dp *di.DependencyProvider, db *sqlx.DB) {
	di.RegisterSingleton(dp, func() *sqlx.DB {
		return db
	})
	di.RegisterCloser(dp, func(ctx context.Context, tag string) (*sqlx.Tx, di.Closer, error) {
		db, err := di.Resolve[*sqlx.DB](ctx, dp)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to resolve database: %w", err)
		}

		tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to start transaction: %w", err)
		}

		return tx, func(err error) error {
			if err == nil || err == context.Canceled {
				txErr := tx.Commit()
				if txErr != nil {
					return txErr
				}
			} else {
				txErr := tx.Rollback()
				if txErr != nil {
					return txErr
				}
			}
			return nil
		}, nil

	})
}
