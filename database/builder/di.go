package builder

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

type contextKey uint8

const (
	txKey contextKey = iota
)

func InitDI(ctx context.Context) error {
	di.Register(func(ctx context.Context, tag string) *sqlx.Tx {
		tx, ok := ctx.Value(txKey).(*sqlx.Tx)
		if ok {
			return tx
		}

		db, ok := di.Resolve[*sqlx.DB](ctx)
		if !ok {
			log.Print("no database register in di")
			return nil
		}

		tx, err := db.BeginTxx(ctx, &sql.TxOptions{
			ReadOnly: strings.ToLower(tag) == "r",
		})
		if err != nil {
			log.Printf("failed to create transaction in di: %v", err)
			return nil
		}
		return tx

	})
	return nil
}

func Transaction(ctx context.Context, readOnly bool, cb func(context.Context) error) error {

	db, ok := di.Resolve[*sqlx.DB](ctx)
	if !ok {
		return fmt.Errorf("no database register in di")
	}

	tx, err := db.BeginTxx(ctx, &sql.TxOptions{
		ReadOnly: readOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	ctx = context.WithValue(ctx, txKey, tx)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	err = cb(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
