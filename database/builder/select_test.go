package builder_test

import (
	"testing"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/test"
)

func TestSelect(t *testing.T) {
	test.QueryTest(t, []test.Case[dialects.QueryBuilder]{
		{
			Name:             "one select",
			Builder:          NewTestBuilder().Select("a"),
			ExpectedSQL:      "SELECT \"a\" FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "two select",
			Builder:          NewTestBuilder().Select("a", "b"),
			ExpectedSQL:      "SELECT \"a\", \"b\" FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "different table",
			Builder:          NewTestBuilder().Select("a.b"),
			ExpectedSQL:      "SELECT \"a\".\"b\" FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "distinct",
			Builder:          NewTestBuilder().Select("a").Distinct(),
			ExpectedSQL:      "SELECT DISTINCT \"a\" FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "subquery",
			Builder:          NewTestBuilder().SelectSubquery(NewTestBuilder().Select("a"), "test"),
			ExpectedSQL:      "SELECT (SELECT \"a\" FROM \"foos\") AS \"test\" FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "all columns",
			Builder:          NewTestBuilder().Select("*"),
			ExpectedSQL:      "SELECT * FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "all columns from table",
			Builder:          NewTestBuilder().Select("foos.*"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\"",
			ExpectedBindings: []any{},
		},
	})
}
