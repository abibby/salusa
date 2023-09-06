package builder_test

import (
	"testing"

	"github.com/abibby/salusa/database/internal/test"
)

func TestOrderBy(t *testing.T) {
	test.QueryTest(t, []test.Case{
		{
			Name:             "one group",
			Builder:          NewTestBuilder().OrderBy("a"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" ORDER BY \"a\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "two groups",
			Builder:          NewTestBuilder().OrderBy("a").OrderBy("b"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" ORDER BY \"a\", \"b\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "different table",
			Builder:          NewTestBuilder().OrderBy("a.b"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" ORDER BY \"a\".\"b\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "descending",
			Builder:          NewTestBuilder().OrderByDesc("a"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" ORDER BY \"a\" DESC",
			ExpectedBindings: []any{},
		},
	})
}
