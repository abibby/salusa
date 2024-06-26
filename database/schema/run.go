package schema

import (
	"context"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type Runner interface {
	Run(ctx context.Context, tx database.DB) error
}

type RunnerFunc func(ctx context.Context, tx database.DB) error

func (f RunnerFunc) Run(ctx context.Context, tx database.DB) error {
	return f(ctx, tx)
}

func Run(f RunnerFunc) Runner {
	return f
}

func runQuery(ctx context.Context, db database.DB, sqler helpers.SQLStringer) error {
	sql, bindings, err := sqler.SQLString(dialects.New())
	if err != nil {
		return err
	}
	_, err = database.Exec(ctx, db, sql, bindings)
	return err
}
