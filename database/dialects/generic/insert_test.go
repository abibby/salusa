package generic_test

import (
	"testing"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/generic"
	"github.com/abibby/salusa/internal/test"
)

func TestGeneric_EncodeInsertQuery(t *testing.T) {
	g := generic.New(&testCore{})
	test.EncoderTest(t, g.EncodeInsertQuery, []test.EncoderTestCase[*dialects.InsertQuery]{
		{
			Name: "single",
			Builder: &dialects.InsertQuery{
				Table:  "foo",
				Values: []map[string]any{{"a": "b"}},
			},
			ExpectedSQL:      "INSERT INTO `foo` (`a`) VALUES (?)",
			ExpectedBindings: []any{"b"},
		},
		{
			Name: "multi",
			Builder: &dialects.InsertQuery{
				Table:  "foo",
				Values: []map[string]any{{"a": "b", "c": 1}, {"a": "d", "c": 2}},
			},
			ExpectedSQL:      "INSERT INTO `foo` (`a`, `c`) VALUES (?, ?), (?, ?)",
			ExpectedBindings: []any{"b", 1, "d", 2},
		},
	})
}
