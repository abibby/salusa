package di_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFill(t *testing.T) {
	t.Run("fill", func(t *testing.T) {
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

		f := &Fillable{}
		err := di.Fill(ctx, f)
		assert.NoError(t, err)
		assert.NotNil(t, f.WithTag)
		assert.Nil(t, f.NoTag)
	})

	t.Run("deep", func(t *testing.T) {
		type Struct struct{ V int }
		type FillableB struct {
			di.Fillable
			Struct *Struct `inject:""`
		}
		type FillableA struct {
			di.Fillable
			B *FillableB `inject:""`
		}
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() *Struct {
			return &Struct{}
		})

		f := &FillableA{}
		err := di.Fill(ctx, f)
		assert.NoError(t, err)
		assert.NotNil(t, f.B)
		assert.NotNil(t, f.B.Struct)
	})

	t.Run("tag", func(t *testing.T) {
		type Struct struct {
			Value string
		}
		type Fillable struct {
			di.Fillable
			WithTag *Struct `inject:"with tag"`
			Empty   *Struct `inject:""`
		}
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.Register(ctx, func(ctx context.Context, tag string) (*Struct, error) {
			return &Struct{
				Value: tag,
			}, nil
		})

		f := &Fillable{}
		err := di.Fill(ctx, f)
		assert.NoError(t, err)
		assert.Equal(t, "with tag", f.WithTag.Value)
		assert.NotNil(t, "", f.Empty.Value)
	})

	t.Run("not registered", func(t *testing.T) {
		type Fillable struct {
			di.Fillable
			Miss *int `inject:""`
		}

		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		f := &Fillable{}
		err := di.Fill(ctx, f)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("non pointer", func(t *testing.T) {
		type Fillable struct {
			di.Fillable
		}

		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		f := Fillable{}
		err := di.Fill(ctx, f)
		assert.ErrorIs(t, err, di.ErrFillParameters)
	})

	t.Run("not struct", func(t *testing.T) {
		type Fillable int

		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		f := Fillable(0)
		err := di.Fill(ctx, &f)
		assert.ErrorIs(t, err, di.ErrFillParameters)
	})

	t.Run("resolve error", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable struct {
			di.Fillable
			S *Struct `inject:""`
		}

		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		resolveErr := fmt.Errorf("error in register")
		di.Register(ctx, func(ctx context.Context, tag string) (*Struct, error) {
			return nil, resolveErr
		})

		f := &Fillable{}
		err := di.Fill(ctx, f)

		assert.Error(t, err)
		assert.ErrorIs(t, err, resolveErr)
	})

	t.Run("auto resolve", func(t *testing.T) {
		type StructA struct{ V int }
		type StructB struct{ V int }
		type Fillable struct {
			di.Fillable
			WithAutoResolve *StructA
			NoTag           *StructB
		}

		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() *StructA {
			return &StructA{}
		})
		di.RegisterSingleton(ctx, func() *StructB {
			return &StructB{}
		})

		f := &Fillable{}
		err := di.Fill(ctx, f, di.AutoResolve[*StructA]())
		assert.NoError(t, err)
		assert.NotNil(t, f.WithAutoResolve)
		assert.Nil(t, f.NoTag)
	})

	t.Run("not fill unfillable type", func(t *testing.T) {
		type Struct struct{ V int }
		type IsFillable struct {
			di.Fillable
			S Struct `inject:""`
		}
		type NotFillable struct {
			S Struct `inject:""`
		}
		type FillableRoot struct {
			IsFillable  *IsFillable  `inject:""`
			NotFillable *NotFillable `inject:""`
		}

		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() *Struct {
			return &Struct{}
		})
		f := &FillableRoot{}
		err := di.Fill(ctx, f)
		assert.ErrorIs(t, err, di.ErrNotRegistered)
	})

	t.Run("deep", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable2 struct {
			di.Fillable
			S *Struct `inject:""`
		}
		type Fillable1 struct {
			di.Fillable
			F *Fillable2 `inject:""`
		}
		type FillableRoot struct {
			F *Fillable1 `inject:""`
		}

		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() *Struct {
			return &Struct{}
		})
		f := &FillableRoot{}
		err := di.Fill(ctx, f)
		assert.NoError(t, err)
		assert.NotNil(t, f.F.F.S)
	})

	t.Run("deep tag", func(t *testing.T) {
		type Struct struct{ V int }
		type Fillable2 struct {
			S *Struct `inject:""`
		}
		type Fillable1 struct {
			F *Fillable2 `inject:"fill"`
		}
		type FillableRoot struct {
			F *Fillable1 `inject:"fill"`
		}

		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() *Struct {
			return &Struct{}
		})
		f := &FillableRoot{}
		err := di.Fill(ctx, f)
		assert.NoError(t, err)
		spew.Dump(f)
		assert.NotNil(t, f.F.F.S)
	})
}
