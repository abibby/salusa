package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/internal/helpers"
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
	di.RegisterWith(ctx, func(ctx context.Context, tag string, deps *userRegisterDeps) (T, error) {
		var zero T
		if deps.Claims == nil {
			return zero, nil
		}
		pkeyValues, err := helpers.PrimaryKeyValue(helpers.NewOf[T]())
		if err != nil {
			return zero, err
		}

		if len(pkeyValues) != 1 {
			return zero, fmt.Errorf("can only inject users with 1 primary key")
		}
		var pkey any
		switch pkeyValues[0].(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			err := json.Unmarshal([]byte(deps.Claims.Subject), &pkey)
			if err != nil {
				return zero, fmt.Errorf("invalid subject must be a number: %w", err)
			}
		default:
			pkey = deps.Claims.Subject
		}

		return builder.From[T]().WithContext(ctx).Find(deps.DB, pkey)
	})
	return nil
}
