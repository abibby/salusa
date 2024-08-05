package di_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestFunc(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type Struct struct{ V int }
		ctx := di.TestDependencyProviderContext()
		di.RegisterSingleton(ctx, func() *Struct {
			return &Struct{V: 10}
		})

		fn := di.PrepareFunc[func(ctx context.Context, i int) int](func(ctx context.Context, i int, s *Struct) int {
			return i + s.V
		})

		result := fn(ctx, 7)

		assert.Equal(t, 17, result)
	})
}
