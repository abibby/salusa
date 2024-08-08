package di_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestDependencyProvider_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		type WithFloat struct {
			Float float64 `inject:""`
		}
		type WithInt struct {
			Int int `inject:""`
		}
		type WithString struct {
			String string `inject:""`
		}
		dp := di.NewDependencyProvider()
		dp.Register(&di.SingletonFactory[float64]{})
		dp.Register(&di.LazySingletonWithFactory[string, *WithFloat]{})
		dp.Register(&di.LazySingletonWithFactory[int, *WithString]{})
		dp.Register(&di.LazySingletonWithFactory[uint, *WithInt]{})

		ctx := context.Background()

		err := dp.Validate(ctx)
		assert.NoError(t, err)
	})

	t.Run("missing dependencies", func(t *testing.T) {
		type WithFloat struct {
			Float float64 `inject:""`
		}
		dp := di.NewDependencyProvider()
		dp.Register(&di.LazySingletonWithFactory[string, *WithFloat]{})

		ctx := context.Background()

		err := dp.Validate(ctx)
		assert.ErrorIs(t, err, di.ErrMissingDependancy)
	})

	t.Run("cycle", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		dp.Register(&di.LazySingletonWithFactory[int, float64]{})
		dp.Register(&di.LazySingletonWithFactory[float64, string]{})
		dp.Register(&di.LazySingletonWithFactory[string, int]{})

		ctx := context.Background()

		err := dp.Validate(ctx)
		assert.ErrorIs(t, err, di.ErrDependancyCycle)
	})

	t.Run("cycle fill", func(t *testing.T) {
		type WithFloat struct {
			Float float64 `inject:""`
		}
		type WithInt struct {
			Int int `inject:""`
		}
		type WithString struct {
			String string `inject:""`
		}
		dp := di.NewDependencyProvider()
		dp.Register(&di.LazySingletonWithFactory[int, *WithFloat]{})
		dp.Register(&di.LazySingletonWithFactory[float64, *WithString]{})
		dp.Register(&di.LazySingletonWithFactory[string, *WithInt]{})

		ctx := context.Background()

		err := dp.Validate(ctx)
		assert.ErrorIs(t, err, di.ErrDependancyCycle)
	})
}
