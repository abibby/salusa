package databasedi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

func Register(dp *di.DependencyProvider, db *sqlx.DB) {
	di.RegisterSingleton(dp, func() *sqlx.DB {
		return db
	})
	di.Register(dp, func(ctx context.Context, tag string) (*sqlx.Tx, error) {
		tx := ctx.Value(txKey)
		if tx == nil {
			tx = &txWrapper{}
		}
		wrapper := tx.(*txWrapper)
		if wrapper.tx == nil {
			db, err := di.Resolve[*sqlx.DB](ctx, dp)
			if errors.Is(err, di.ErrNotRegistered) {
				return nil, fmt.Errorf("the database is not registered in di")
			} else if err != nil {
				return nil, fmt.Errorf("failed to open database: %w", err)
			}
			tx, err := db.BeginTxx(ctx, &sql.TxOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to start transaction: %w", err)
			}
			wrapper.tx = tx
		}
		return wrapper.tx, nil

	})
}
