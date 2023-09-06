package helpers_test

import (
	"reflect"
	"testing"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	type Foo struct {
		ID  int
		Foo string `db:"foo,autoincrement,primary,type:date"`
	}

	rt := reflect.TypeOf(Foo{})
	assert.Equal(
		t,
		&helpers.Tag{
			Name:          "ID",
			Primary:       false,
			AutoIncrement: false,
			Readonly:      false,
			Index:         false,
		},
		helpers.DBTag(rt.Field(0)),
	)
	assert.Equal(
		t,
		&helpers.Tag{
			Name:          "foo",
			Primary:       true,
			AutoIncrement: true,
			Readonly:      false,
			Index:         false,
			Type:          dialects.DataType("date"),
		},
		helpers.DBTag(rt.Field(1)),
	)
}
