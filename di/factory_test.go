package di_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestNewFactoryFuncWith(t *testing.T) {
	t.Run("with", func(t *testing.T) {
		type Parent struct{ V int }
		type Child struct{ Parent *Parent }
		ctx := context.Background()
		dp := di.NewDependencyProvider()

		p := &Parent{V: 10}
		dp.Register(di.NewSingletonFactory(p))

		c, err := di.NewFactoryFuncWith(func(ctx context.Context, tag string, parent *Parent) (*Child, error) {
			return &Child{Parent: parent}, nil
		}).Build(ctx, dp, "")

		assert.NoError(t, err)
		assert.Equal(t, p, c.(*Child).Parent)
	})

}
