package auth

import (
	"context"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

type userRegisterDeps struct {
	DB     *sqlx.DB `inject:""`
	Claims *Claims  `inject:""`
}

func Register[T User](ctx context.Context) error {
	di.Register(ctx, func(ctx context.Context, tag string) (*Claims, error) {
		c, _ := GetClaimsCtx(ctx)
		return c, nil
	})
	di.Register(ctx, func(ctx context.Context, tag string) (T, error) {
		deps := &userRegisterDeps{}
		err := di.Fill(ctx, deps)
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
	return nil
}
