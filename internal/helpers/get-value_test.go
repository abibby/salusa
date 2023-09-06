package helpers_test

import (
	"testing"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestGetValue(t *testing.T) {
	type A struct {
		Foo string
	}
	type B struct {
		A
	}

	type C struct {
		Foo *A `db:"foo"`
	}
	var nilA *A
	var nilC *C

	testCases := []struct {
		Name          string
		Value         any
		Field         string
		ExpectedValue any
		ExpectedOk    bool
	}{
		{"heap find", A{Foo: "bar"}, "Foo", "bar", true},
		{"heap miss", A{}, "Bar", nil, false},
		{"pointer find", &A{Foo: "bar"}, "Foo", "bar", true},
		{"pointer miss", &A{}, "Bar", nil, false},
		{"nil pointer find", nilA, "Foo", nil, false},
		{"nil pointer miss", nilA, "Bar", nil, false},
		{"anonymise find", &B{A: A{Foo: "bar"}}, "Foo", "bar", true},
		{"anonymise miss", &B{}, "Bar", nil, false},
		{"nil tag find", nilC, "foo", nil, false},
		{"nil tag miss Foo", nilC, "Foo", nil, false},
		{"nil tag miss", nilC, "bar", nil, false},
		{"nil value", C{}, "foo", (*A)(nil), true},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			v, ok := helpers.GetValue(tc.Value, tc.Field)
			assert.Equal(t, tc.ExpectedOk, ok)
			assert.Equal(t, tc.ExpectedValue, v)
		})
	}
}
