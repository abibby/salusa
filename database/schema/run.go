package schema

import (
	"context"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/internal/helpers"
)

type Runner interface {
	Run(ctx context.Context, tx helpers.QueryExecer) error
}

type RunnerFunc func(ctx context.Context, tx helpers.QueryExecer) error

func (f RunnerFunc) Run(ctx context.Context, tx helpers.QueryExecer) error {
	return f(ctx, tx)
}

func Run(f RunnerFunc) Runner {
	return f
}

func runQuery(ctx context.Context, tx helpers.QueryExecer, sqler helpers.ToSQLer) error {
	sql, bindings, err := sqler.ToSQL(dialects.New())
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, sql, bindings...)
	return err
}
