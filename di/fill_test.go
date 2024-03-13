package di_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestFill(t *testing.T) {
	t.Run("fill", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			WithTag *Struct `inject:""`
			NoTag   *Struct
		}

		dp := di.NewDependencyProvider()
		di.RegisterSingleton(dp, func() *Struct {
			return &Struct{}
		})

		f := &Fillable{}
		err := dp.Fill(context.Background(), f)
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
		dp := di.NewDependencyProvider()
		di.RegisterSingleton(dp, func() *Struct {
			return &Struct{}
		})

		f := &FillableA{}
		err := dp.Fill(context.Background(), f)
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
		dp := di.NewDependencyProvider()
		di.Register(dp, func(ctx context.Context, tag string) (*Struct, error) {
			return &Struct{
				Value: tag,
			}, nil
		})

		f := &Fillable{}
		err := dp.Fill(context.Background(), f)
		assert.NoError(t, err)
		assert.Equal(t, "with tag", f.WithTag.Value)
		assert.NotNil(t, "", f.Empty.Value)
	})

	t.Run("not registered", func(t *testing.T) {
		type Fillable struct {
			Miss *int `inject:""`
		}

		dp := di.NewDependencyProvider()
		f := &Fillable{}
		err := dp.Fill(context.Background(), f)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("non pointer", func(t *testing.T) {
		type Fillable struct{}

		dp := di.NewDependencyProvider()
		f := Fillable{}
		err := dp.Fill(context.Background(), f)
		assert.ErrorIs(t, err, di.ErrFillParameters)
	})

	t.Run("not struct", func(t *testing.T) {
		type Fillable int

		dp := di.NewDependencyProvider()
		f := Fillable(0)
		err := dp.Fill(context.Background(), &f)
		assert.ErrorIs(t, err, di.ErrFillParameters)
	})

	t.Run("resolve error", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			S *Struct `inject:""`
		}

		dp := di.NewDependencyProvider()
		resolveErr := fmt.Errorf("error in register")
		di.Register(dp, func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		f := &Fillable{}
		err := dp.Fill(context.Background(), f)

		assert.Error(t, err)
		assert.ErrorIs(t, err, resolveErr)
	})

	t.Run("auto resolve", func(t *testing.T) {
		type StructA struct{ V int }
		type StructB struct{ V int }
		type Fillable struct {
			WithAutoResolve *StructA
			NoTag           *StructB
		}

		dp := di.NewDependencyProvider()
		di.RegisterSingleton(dp, func() *StructA {
			return &StructA{}
		})
		di.RegisterSingleton(dp, func() *StructB {
			return &StructB{}
		})

		f := &Fillable{}
		err := dp.Fill(context.Background(), f, di.AutoResolve[*StructA]())
		assert.NoError(t, err)
		assert.NotNil(t, f.WithAutoResolve)
		assert.Nil(t, f.NoTag)
	})
}
