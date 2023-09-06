package builder_test

import (
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/internal/test"
)

func TestHaving(t *testing.T) {
	test.QueryTest(t, []test.Case{
		{
			Name:             "one where",
			Builder:          NewTestBuilder().Having("a", "=", "b"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" HAVING \"a\" = ?",
			ExpectedBindings: []any{"b"},
		},
		{
			Name:             "2 wheres",
			Builder:          NewTestBuilder().Having("a", "=", "b").Having("c", "=", "d"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" HAVING \"a\" = ? AND \"c\" = ?",
			ExpectedBindings: []any{"b", "d"},
		},
		{
			Name:             "specified table",
			Builder:          NewTestBuilder().Having("foo.a", "=", "b"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" HAVING \"foo\".\"a\" = ?",
			ExpectedBindings: []any{"b"},
		},
		{
			Name:             "or where",
			Builder:          NewTestBuilder().Having("a", "=", "b").OrHaving("c", "=", "d"),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" HAVING \"a\" = ? OR \"c\" = ?",
			ExpectedBindings: []any{"b", "d"},
		},
		{
			Name: "and group",
			Builder: NewTestBuilder().HavingAnd(func(wl *builder.Conditions) {
				wl.Where("a", "=", "a").OrWhere("b", "=", "b")
			}).HavingAnd(func(wl *builder.Conditions) {
				wl.Where("c", "=", "c").OrWhere("d", "=", "d")
			}),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" HAVING (\"a\" = ? OR \"b\" = ?) AND (\"c\" = ? OR \"d\" = ?)",
			ExpectedBindings: []any{"a", "b", "c", "d"},
		},
		{
			Name: "or group",
			Builder: NewTestBuilder().HavingOr(func(wl *builder.Conditions) {
				wl.Where("a", "=", "a").Where("b", "=", "b")
			}).HavingOr(func(wl *builder.Conditions) {
				wl.Where("c", "=", "c").Where("d", "=", "d")
			}),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" HAVING (\"a\" = ? AND \"b\" = ?) OR (\"c\" = ? AND \"d\" = ?)",
			ExpectedBindings: []any{"a", "b", "c", "d"},
		},
		{
			Name:             "subquery",
			Builder:          NewTestBuilder().Having("a", "=", NewTestBuilder().Select("a").Having("id", "=", 1)),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" HAVING \"a\" = (SELECT \"a\" FROM \"foos\" HAVING \"id\" = ?)",
			ExpectedBindings: []any{1},
		},
	})
}
