package di_test

import (
	"context"
	"os"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	di.SetDefaultProvider(di.NewDependencyProvider())

	code := m.Run()

	os.Exit(code)
}

func TestRegister(t *testing.T) {
	t.Run("singlton", func(t *testing.T) {
		type Struct struct{}
		di.RegisterSingleton(func() *Struct {
			return &Struct{}
		})
		ctx := context.Background()
		a, aOk := di.Resolve[*Struct](ctx)
		b, bOk := di.Resolve[*Struct](ctx)
		assert.NotNil(t, a)
		assert.True(t, aOk)
		assert.NotNil(t, b)
		assert.True(t, bOk)
		assert.Same(t, a, b)
	})
	t.Run("interface", func(t *testing.T) {
		type Interface interface{}
		type Struct struct{}
		di.Register(func(ctx context.Context, tag string) Interface {
			return &Struct{}
		})

		s, ok := di.Resolve[Interface](context.Background())
		assert.NotNil(t, s)
		assert.True(t, ok)

		_, ok = s.(*Struct)
		assert.True(t, ok)
	})

	t.Run("non singleton", func(t *testing.T) {
		type Struct struct {
			A int
		}
		i := 0
		di.Register(func(ctx context.Context, tag string) *Struct {
			i++
			return &Struct{
				A: i,
			}
		})

		ctx := context.Background()
		a, _ := di.Resolve[*Struct](ctx)
		b, _ := di.Resolve[*Struct](ctx)

		assert.Equal(t, 1, a.A)
		assert.Equal(t, 2, b.A)
	})

	t.Run("not registered", func(t *testing.T) {
		type Struct struct{}

		v, ok := di.Resolve[*Struct](context.Background())
		assert.Nil(t, v)
		assert.False(t, ok)
	})

	t.Run("same name", func(t *testing.T) {
		{
			type Struct struct{}
			di.RegisterSingleton(func() *Struct {
				return &Struct{}
			})
		}
		{
			type Struct struct{}
			v, ok := di.Resolve[*Struct](context.Background())
			assert.Nil(t, v)
			assert.False(t, ok)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		assert.PanicsWithValue(t, di.ErrInvalidDependencyFactory, func() {
			dp := di.NewDependencyProvider()
			dp.Register(func() {})
		})
	})
}

func TestFill(t *testing.T) {
	t.Run("fill", func(t *testing.T) {
		type Struct struct{}
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
		type Struct struct{}
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
		di.Register(func(ctx context.Context, tag string) *Struct {
			return &Struct{
				Value: tag,
			}
		})

		f := &Fillable{}
		err := di.Fill(context.Background(), f)
		assert.NoError(t, err)
		assert.Equal(t, "with tag", f.WithTag.Value)
		assert.NotNil(t, "", f.Empty.Value)
	})
}
