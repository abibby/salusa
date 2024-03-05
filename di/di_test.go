package di_test

import (
	"os"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	di.DefaultProvider = di.NewDependamcyProvider()

	code := m.Run()

	os.Exit(code)
}

func TestRegister(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type Struct struct{}
		ogS := &Struct{}
		di.Register[*Struct]().Singlton(ogS)

		s, ok := di.Resolve[*Struct]()
		assert.NotNil(t, s)
		assert.True(t, ok)
		assert.Same(t, ogS, s)
	})
	t.Run("interface", func(t *testing.T) {
		type Interface interface{}
		type Struct struct{}
		ogS := &Struct{}
		di.Register[Interface]().Singlton(ogS)

		s, ok := di.Resolve[Interface]()
		assert.NotNil(t, s)
		assert.True(t, ok)
		assert.Same(t, ogS, s)
	})

	t.Run("non singleton", func(t *testing.T) {
		type Struct struct {
			A int
		}
		i := 0
		di.Register[*Struct]().Factory(func() *Struct {
			i++
			return &Struct{
				A: i,
			}
		})

		a, _ := di.Resolve[*Struct]()
		b, _ := di.Resolve[*Struct]()

		assert.Equal(t, 1, a.A)
		assert.Equal(t, 2, b.A)
	})

	t.Run("not registered", func(t *testing.T) {
		type Struct struct{}

		v, ok := di.Resolve[*Struct]()
		assert.Nil(t, v)
		assert.False(t, ok)
	})

	t.Run("same name", func(t *testing.T) {
		{
			type Struct struct{}
			di.Register[*Struct]().Singlton(&Struct{})
		}
		{
			type Struct struct{}
			v, ok := di.Resolve[*Struct]()
			assert.Nil(t, v)
			assert.False(t, ok)
		}
	})
}

func TestFill(t *testing.T) {
	t.Run("fill", func(t *testing.T) {
		type Struct struct{}
		type Fillable struct {
			S *Struct `di:"inject"` // maybe struct tags
		}
		ogS := &Struct{}
		di.Register[*Struct]().Singlton(ogS)

		f := &Fillable{}
		err := di.Fill(f)
		assert.NoError(t, err)
		assert.Same(t, ogS, f.S)
	})
}
