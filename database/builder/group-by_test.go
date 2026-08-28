package builder_test

import (
	"testing"

	"abibby.com/salusa/database/dialects"
	"abibby.com/salusa/internal/test"
)

func TestGroupBy(t *testing.T) {
	test.QueryTest(t, []test.Case[dialects.QueryBuilder]{
		{
			Name:             "one group",
			Builder:          NewTestBuilder().GroupBy("a"),
			ExpectedSQLite:   "SELECT \"foos\".* FROM \"foos\" GROUP BY \"a\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "two groups",
			Builder:          NewTestBuilder().GroupBy("a", "b"),
			ExpectedSQLite:   "SELECT \"foos\".* FROM \"foos\" GROUP BY \"a\", \"b\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "different table",
			Builder:          NewTestBuilder().GroupBy("a.b"),
			ExpectedSQLite:   "SELECT \"foos\".* FROM \"foos\" GROUP BY \"a\".\"b\"",
			ExpectedBindings: []any{},
		},
	})
}
