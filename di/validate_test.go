package di_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestDependencyProvider_Validate(t *testing.T) {
	t.Run("cycle", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		dp.Register(&di.LazySingletonWithFactory[int, float64]{})
		dp.Register(&di.LazySingletonWithFactory[float64, string]{})
		dp.Register(&di.LazySingletonWithFactory[string, int]{})

		ctx := context.Background()

		err := dp.Validate(ctx)
		// assert.NoError(t, err)
		assert.ErrorIs(t, err, di.ErrDependancyCycle)
	})
}
