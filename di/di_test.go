package di_test

import (
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	type Interface interface{}
	type Struct struct{}
	type Fillable struct {
		S *Struct `di:"inject"` // maybe struct tags
	}

	di.Register[*Struct]().Singlton(&Struct{})
	di.Register[Interface]().Singlton(&Struct{})
	di.Register[Interface]().Factory(func() Interface {
		return &Struct{}
	})

	f := &Fillable{}

	di.Fill(f)

	assert.NotNil(t, f.S)

	s := di.Resolve[*Struct]()
	assert.NotNil(t, s)
}
