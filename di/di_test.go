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
		dp := di.NewDependencyProvider()
		di.RegisterSingleton(dp, func() *Struct {
			return &Struct{}
		})
		ctx := context.Background()
		a, aErr := di.Resolve[*Struct](ctx, dp)
		b, bErr := di.Resolve[*Struct](ctx, dp)
		assert.NotNil(t, a)
		assert.NoError(t, aErr)
		assert.NotNil(t, b)
		assert.NoError(t, bErr)
		assert.Same(t, a, b)
	})

	t.Run("interface", func(t *testing.T) {
		type Interface interface{}
		dp := di.NewDependencyProvider()
		type Struct struct{ V int }
		di.Register(dp, func(ctx context.Context, tag string) (Interface, error) {
			return &Struct{}, nil
		})

		s, err := di.Resolve[Interface](context.Background(), dp)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		_, ok := s.(*Struct)
		assert.True(t, ok)
	})

	t.Run("non singleton", func(t *testing.T) {
		type Struct struct {
			A int
		}
		dp := di.NewDependencyProvider()
		i := 0
		di.Register(dp, func(ctx context.Context, tag string) (*Struct, error) {
			i++
			return &Struct{
				A: i,
			}, nil
		})

		ctx := context.Background()
		a, _ := di.Resolve[*Struct](ctx, dp)
		b, _ := di.Resolve[*Struct](ctx, dp)

		assert.Equal(t, 1, a.A)
		assert.Equal(t, 2, b.A)
	})

	t.Run("not registered", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		v, err := di.Resolve[int](context.Background(), dp)
		assert.Zero(t, v)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("same name", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		{
			type Struct int
			di.RegisterSingleton(dp, func() Struct {
				return 123
			})
		}
		{
			type Struct int
			v, err := di.Resolve[Struct](context.Background(), dp)
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
		dp := di.NewDependencyProvider()

		di.RegisterSingleton(dp, func() *Struct {
			return &Struct{}
		})

		f, err := di.Resolve[*Fillable](context.Background(), dp)
		assert.NoError(t, err)
		assert.NotNil(t, f.WithTag)
		assert.Nil(t, f.NoTag)
	})

	t.Run("basic type", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		di.RegisterSingleton(dp, func() int {
			return 123
		})

		i, err := di.Resolve[int](context.Background(), dp)
		assert.NoError(t, err)
		assert.Equal(t, 123, i)
	})
}

func TestRegisterLazySingleton(t *testing.T) {
	t.Run("is lazy", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		type Struct struct{ V int }
		runs := 0
		di.RegisterLazySingleton(dp, func() *Struct {
			runs++
			return &Struct{}
		})

		ctx := context.Background()
		assert.Equal(t, 0, runs)
		a, aErr := di.Resolve[*Struct](ctx, dp)
		assert.Equal(t, 1, runs)
		b, bErr := di.Resolve[*Struct](ctx, dp)
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
		expectedContext := context.WithValue(context.Background(), "foo", "bar")

		dp := di.NewDependencyProvider()
		ctx, err := di.Resolve[context.Context](expectedContext, dp)

		assert.NoError(t, err)
		assert.Same(t, expectedContext, ctx)
	})

	t.Run("self", func(t *testing.T) {

		expectedDP := di.NewDependencyProvider()
		dp, err := di.Resolve[*di.DependencyProvider](context.Background(), expectedDP)

		assert.NoError(t, err)
		assert.Same(t, expectedDP, dp)
	})

	t.Run("error", func(t *testing.T) {
		type Struct struct{ V int }
		dp := di.NewDependencyProvider()
		resolveErr := fmt.Errorf("resolve error")
		di.Register(dp, func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		ctx := context.Background()
		v, err := di.Resolve[*Struct](ctx, dp)

		assert.Same(t, resolveErr, err)
		assert.Zero(t, v)
	})

	t.Run("error in fill", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			di.Fillable
			S *Struct `inject:""`
		}
		dp := di.NewDependencyProvider()
		resolveErr := fmt.Errorf("resolve error")
		di.Register(dp, func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		ctx := context.Background()
		v, err := di.Resolve[*Fillable](ctx, dp)

		assert.ErrorIs(t, err, resolveErr)
		assert.Zero(t, v)
	})
}
