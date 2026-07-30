package generic_test

import (
	"testing"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/generic"
	"github.com/abibby/salusa/internal/test"
)

func TestGeneric_EncodeConditions(t *testing.T) {

	g := generic.New(&testCore{})
	test.EncoderTest(t, g.EncodeConditions, []test.EncoderTestCase[[]dialects.Condition]{
		{
			Name:             "empty",
			Builder:          []dialects.Condition{},
			ExpectedSQL:      "",
			ExpectedBindings: []any{},
		},
		{
			Name: "single simple",
			Builder: []dialects.Condition{
				{
					Column:   dialects.Column{Column: "foo"},
					Operator: "=",
					Value:    "bar",
				},
			},
			ExpectedSQL:      "foo = ?",
			ExpectedBindings: []any{"bar"},
		},
	})
}
