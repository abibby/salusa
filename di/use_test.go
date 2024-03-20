package di_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestUse(t *testing.T) {
	t.Run("use", func(t *testing.T) {
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

		err := di.Use(context.Background(), dp, func(f *Fillable) error {
			assert.NotNil(t, f.WithTag)
			assert.Nil(t, f.NoTag)
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("closer", func(t *testing.T) {
		type Struct struct{ Closed bool }
		type Fillable struct {
			Struct *Struct `inject:""`
		}

		dp := di.NewDependencyProvider()
		di.RegisterCloser(dp, func(ctx context.Context, tag string) (*Struct, di.Closer, error) {
			s := &Struct{
				Closed: false,
			}
			return s, func(err error) error {
				s.Closed = true
				return nil
			}, nil
		})
		var s *Struct
		err := di.Use(context.Background(), dp, func(f *Fillable) error {
			assert.NotNil(t, f.Struct)
			s = f.Struct

			return nil
		})
		assert.NoError(t, err)
		assert.True(t, s.Closed)
	})
}
