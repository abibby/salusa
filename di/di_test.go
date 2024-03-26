package di_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Run("singlton", func(t *testing.T) {
		type Struct struct{ V int }
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() *Struct {
			return &Struct{}
		})
		a, aErr := di.Resolve[*Struct](ctx)
		b, bErr := di.Resolve[*Struct](ctx)
		assert.NotNil(t, a)
		assert.NoError(t, aErr)
		assert.NotNil(t, b)
		assert.NoError(t, bErr)
		assert.Same(t, a, b)
	})

	t.Run("interface", func(t *testing.T) {
		type Interface interface{}
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		type Struct struct{ V int }
		di.Register(ctx, func(ctx context.Context, tag string) (Interface, error) {
			return &Struct{}, nil
		})

		s, err := di.Resolve[Interface](ctx)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		_, ok := s.(*Struct)
		assert.True(t, ok)
	})

	t.Run("non singleton", func(t *testing.T) {
		type Struct struct {
			A int
		}
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		i := 0
		di.Register(ctx, func(ctx context.Context, tag string) (*Struct, error) {
			i++
			return &Struct{
				A: i,
			}, nil
		})

		a, _ := di.Resolve[*Struct](ctx)
		b, _ := di.Resolve[*Struct](ctx)

		assert.Equal(t, 1, a.A)
		assert.Equal(t, 2, b.A)
	})

	t.Run("not registered", func(t *testing.T) {
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		v, err := di.Resolve[int](ctx)
		assert.Zero(t, v)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("same name", func(t *testing.T) {
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		{
			type Struct int
			di.RegisterSingleton(ctx, func() Struct {
				return 123
			})
		}
		{
			type Struct int
			v, err := di.Resolve[Struct](ctx)
			assert.Zero(t, v)
			assert.ErrorIs(t, err, di.ErrNotRegistered)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		assert.PanicsWithValue(t, di.ErrInvalidDependencyFactory, func() {
			dp := di.NewDependencyProvider()
			dp.Register(func() {})
		})
	})

	t.Run("resolve fillable struct", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			di.Fillable
			WithTag *Struct `inject:""`
			NoTag   *Struct
		}
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)

		di.RegisterSingleton(ctx, func() *Struct {
			return &Struct{}
		})

		f, err := di.Resolve[*Fillable](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, f.WithTag)
		assert.Nil(t, f.NoTag)
	})

	t.Run("basic type", func(t *testing.T) {
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() int {
			return 123
		})

		i, err := di.Resolve[int](ctx)
		assert.NoError(t, err)
		assert.Equal(t, 123, i)
	})
}

func TestRegisterLazySingleton(t *testing.T) {
	t.Run("is lazy", func(t *testing.T) {
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		type Struct struct{ V int }
		runs := 0
		di.RegisterLazySingleton(ctx, func() *Struct {
			runs++
			return &Struct{}
		})

		assert.Equal(t, 0, runs)
		a, aErr := di.Resolve[*Struct](ctx)
		assert.Equal(t, 1, runs)
		b, bErr := di.Resolve[*Struct](ctx)
		assert.Equal(t, 1, runs)

		assert.NotNil(t, a)
		assert.NoError(t, aErr)
		assert.NotNil(t, b)
		assert.NoError(t, bErr)
		assert.Same(t, a, b)
	})
}

func TestResolve(t *testing.T) {
	t.Run("context", func(t *testing.T) {

		expectedContext := di.ContextWithDependencyProvider(
			context.WithValue(context.Background(), "foo", "bar"),
			di.NewDependencyProvider(),
		)
		ctx, err := di.Resolve[context.Context](expectedContext)

		assert.NoError(t, err)
		assert.Same(t, expectedContext, ctx)
	})

	t.Run("self", func(t *testing.T) {
		expectedDP := di.NewDependencyProvider()
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			expectedDP,
		)
		dp, err := di.Resolve[*di.DependencyProvider](ctx)

		assert.NoError(t, err)
		assert.Same(t, expectedDP, dp)
	})

	t.Run("error", func(t *testing.T) {
		type Struct struct{ V int }
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		resolveErr := fmt.Errorf("resolve error")
		di.Register(ctx, func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		v, err := di.Resolve[*Struct](ctx)

		assert.Same(t, resolveErr, err)
		assert.Zero(t, v)
	})

	t.Run("error in fill", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			di.Fillable
			S *Struct `inject:""`
		}
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		resolveErr := fmt.Errorf("resolve error")
		di.Register(ctx, func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		v, err := di.Resolve[*Fillable](ctx)

		assert.ErrorIs(t, err, resolveErr)
		assert.Zero(t, v)
	})
}
