package di_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	di.TestMain(m)
}

func TestRegister(t *testing.T) {
	t.Run("singlton", func(t *testing.T) {
		type Struct struct{ V int }
		di.RegisterSingleton(func() *Struct {
			return &Struct{}
		})
		ctx := context.Background()
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
		type Struct struct{ V int }
		di.Register(func(ctx context.Context, tag string) (Interface, error) {
			return &Struct{}, nil
		})

		s, err := di.Resolve[Interface](context.Background())
		assert.NotNil(t, s)
		assert.NoError(t, err)

		_, ok := s.(*Struct)
		assert.True(t, ok)
	})

	t.Run("non singleton", func(t *testing.T) {
		type Struct struct {
			A int
		}
		i := 0
		di.Register(func(ctx context.Context, tag string) (*Struct, error) {
			i++
			return &Struct{
				A: i,
			}, nil
		})

		ctx := context.Background()
		a, _ := di.Resolve[*Struct](ctx)
		b, _ := di.Resolve[*Struct](ctx)

		assert.Equal(t, 1, a.A)
		assert.Equal(t, 2, b.A)
	})

	t.Run("not registered", func(t *testing.T) {
		v, err := di.Resolve[int](context.Background())
		assert.Zero(t, v)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("same name", func(t *testing.T) {
		{
			type Struct int
			di.RegisterSingleton(func() Struct {
				return 123
			})
		}
		{
			type Struct int
			v, err := di.Resolve[Struct](context.Background())
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

	t.Run("struct", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			WithTag *Struct `inject:""`
			NoTag   *Struct
		}

		di.RegisterSingleton(func() *Struct {
			return &Struct{}
		})

		f, err := di.Resolve[*Fillable](context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, f.WithTag)
		assert.Nil(t, f.NoTag)
	})

	t.Run("basic type", func(t *testing.T) {
		di.RegisterSingleton(func() int {
			return 123
		})

		i, err := di.Resolve[int](context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 123, i)
	})
}

func TestRegisterLazySingleton(t *testing.T) {
	t.Run("is lazy", func(t *testing.T) {
		type Struct struct{ V int }
		runs := 0
		di.RegisterLazySingleton(func() *Struct {
			runs++
			return &Struct{}
		})

		ctx := context.Background()
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

func TestFill(t *testing.T) {
	t.Run("fill", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			WithTag *Struct `inject:""`
			NoTag   *Struct
		}

		di.RegisterSingleton(func() *Struct {
			return &Struct{}
		})

		f := &Fillable{}
		err := di.Fill(context.Background(), f)
		assert.NoError(t, err)
		assert.NotNil(t, f.WithTag)
		assert.Nil(t, f.NoTag)
	})

	t.Run("deep", func(t *testing.T) {
		type Struct struct{ V int }
		type FillableB struct {
			Struct *Struct `inject:""`
		}
		type FillableA struct {
			B *FillableB `inject:""`
		}
		di.RegisterSingleton(func() *Struct {
			return &Struct{}
		})

		f := &FillableA{}
		err := di.Fill(context.Background(), f)
		assert.NoError(t, err)
		assert.NotNil(t, f.B)
		assert.NotNil(t, f.B.Struct)
	})

	t.Run("tag", func(t *testing.T) {
		type Struct struct {
			Value string
		}
		type Fillable struct {
			WithTag *Struct `inject:"with tag"`
			Empty   *Struct `inject:""`
		}
		di.Register(func(ctx context.Context, tag string) (*Struct, error) {
			return &Struct{
				Value: tag,
			}, nil
		})

		f := &Fillable{}
		err := di.Fill(context.Background(), f)
		assert.NoError(t, err)
		assert.Equal(t, "with tag", f.WithTag.Value)
		assert.NotNil(t, "", f.Empty.Value)
	})

	t.Run("not registered", func(t *testing.T) {
		type Fillable struct {
			Miss *int `inject:""`
		}

		f := &Fillable{}
		err := di.Fill(context.Background(), f)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("non pointer", func(t *testing.T) {
		type Fillable struct{}

		f := Fillable{}
		err := di.Fill(context.Background(), f)
		assert.ErrorIs(t, err, di.ErrFillParameters)
	})

	t.Run("not struct", func(t *testing.T) {
		type Fillable int

		f := Fillable(0)
		err := di.Fill(context.Background(), &f)
		assert.ErrorIs(t, err, di.ErrFillParameters)
	})

	t.Run("resolve error", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			S *Struct `inject:""`
		}

		resolveErr := fmt.Errorf("error in register")
		di.Register(func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		f := &Fillable{}
		err := di.Fill(context.Background(), f)

		assert.Error(t, err)
		assert.ErrorIs(t, err, resolveErr)
	})
}

func TestResolve(t *testing.T) {
	t.Run("context", func(t *testing.T) {
		expectedContext := context.WithValue(context.Background(), "foo", "bar")

		ctx, err := di.Resolve[context.Context](expectedContext)

		assert.NoError(t, err)
		assert.Same(t, expectedContext, ctx)
	})

	t.Run("error", func(t *testing.T) {
		type Struct struct{ V int }
		resolveErr := fmt.Errorf("resolve error")
		di.Register(func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		ctx := context.Background()
		v, err := di.Resolve[*Struct](ctx)

		assert.Same(t, resolveErr, err)
		assert.Zero(t, v)
	})

	t.Run("error in fill", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			S *Struct `inject:""`
		}
		resolveErr := fmt.Errorf("resolve error")
		di.Register(func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		ctx := context.Background()
		v, err := di.Resolve[*Fillable](ctx)

		assert.ErrorIs(t, err, resolveErr)
		assert.Zero(t, v)
	})
}
