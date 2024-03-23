package auth

import (
	"context"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

type userRegisterDeps struct {
	DB     *sqlx.DB
	Claims *Claims
}

func Register[T User](dp *di.DependencyProvider) {
	di.Register(dp, func(ctx context.Context, tag string) (*Claims, error) {
		c, _ := GetClaimsCtx(ctx)
		return c, nil
	})
	di.Register(dp, func(ctx context.Context, tag string) (T, error) {
		deps := &userRegisterDeps{}
		err := dp.Fill(ctx, deps)
		if err != nil {
			var zero T
			return zero, err
		}
		if deps.Claims == nil {
			var zero T
			return zero, nil
		}

		return builder.From[T]().WithContext(ctx).Find(deps.DB, deps.Claims.Subject)
	})
}
