package schema

import (
	"context"

	"abibby.com/salusa/database"
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
