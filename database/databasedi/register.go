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
	di.Register(dp, func(ctx context.Context, tag string) (*sqlx.Tx, error) {
		db, err := di.Resolve[*sqlx.DB](ctx, dp)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve database: %w", err)
		}

		tx, err := db.BeginTxx(context.WithoutCancel(ctx), &sql.TxOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to start transaction: %w", err)
		}

		go func() {
			<-ctx.Done()

			err := context.Cause(ctx)
			if err == nil || err == context.Canceled {
				txErr := tx.Commit()
				if txErr != nil {
					panic(txErr)
				}
			} else {
				txErr := tx.Rollback()
				if txErr != nil {
					panic(txErr)
				}
			}
		}()
		return tx, nil

	})
}
