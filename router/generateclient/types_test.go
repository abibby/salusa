package generateclient

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestToTsType(t *testing.T) {

	testCases := []struct {
		Name     string
		Type     reflect.Type
		Expected string
	}{
		{
			Name:     "string",
			Type:     helpers.GetType[string](),
			Expected: "string",
		},
		{
			Name:     "int",
			Type:     helpers.GetType[int](),
			Expected: "number",
		},
		{
			Name:     "int32",
			Type:     helpers.GetType[int32](),
			Expected: "number",
		},
		{
			Name:     "float32",
			Type:     helpers.GetType[float32](),
			Expected: "number",
		},
		{
			Name:     "pointer",
			Type:     helpers.GetType[*int](),
			Expected: "number | null",
		},
		{
			Name:     "slice",
			Type:     helpers.GetType[[]int](),
			Expected: "Array<number>",
		},
		{
			Name: "struct",
			Type: helpers.GetType[struct {
				Foo string `json:"foo"`
			}](),
			Expected: "{foo:string}",
		},
		{
			Name: "nested struct",
			Type: helpers.GetType[struct {
				Foo struct {
					Bar string `json:"bar"`
				} `json:"foo"`
			}](),
			Expected: "{foo:{bar:string}}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			b := &bytes.Buffer{}
			err := toTsType(tc.Type, b)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, b.String())
		})
	}
}
