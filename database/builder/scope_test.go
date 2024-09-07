package builder_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model/mixins"
	"github.com/abibby/salusa/internal/test"
)

type ScopeFoo struct {
	test.Foo
}

func (f *ScopeFoo) Scopes() []*builder.Scope {
	return []*builder.Scope{
		mixins.SoftDeleteScope,
	}
}

type ScopeBar struct {
	test.Bar
	ScopeFoo *builder.BelongsTo[*ScopeFoo] `db:"-" json:"foo"`
}

func TestScope(t *testing.T) {
	scopeA := &builder.Scope{
		Name: "with-a",
		Query: func(b *builder.Builder) *builder.Builder {
			return b.Where("a", "=", "b")
		},
	}
	scopeCtx := &builder.Scope{
		Name: "ctx",
		Query: func(b *builder.Builder) *builder.Builder {
			foo := b.Context().Value("foo")
			return b.Where("a", "=", foo)
		},
	}
	test.QueryTest(t, []test.Case{
		{
			Name:             "scope",
			Builder:          NewTestBuilder().WithScope(scopeA),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" WHERE \"a\" = ?",
			ExpectedBindings: []any{"b"},
		},
		{
			Name:             "without scope",
			Builder:          NewTestBuilder().WithScope(scopeA).WithoutScope(scopeA),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name:             "global scope",
			Builder:          builder.From[*ScopeFoo](),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" WHERE \"foos\".\"deleted_at\" IS NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "without global scope",
			Builder:          builder.From[*ScopeFoo]().WithoutGlobalScope(mixins.SoftDeleteScope),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\"",
			ExpectedBindings: []any{},
		},
		{
			Name: "global scope whereHas",
			Builder: builder.From[*ScopeBar]().WhereHas("ScopeFoo", func(q *builder.Builder) *builder.Builder {
				return q
			}),
			ExpectedSQL:      `SELECT "bars".* FROM "bars" WHERE EXISTS (SELECT "foos".* FROM "foos" WHERE "id" = "bars"."foo_id" AND "foos"."deleted_at" IS NULL)`,
			ExpectedBindings: []any{},
		},
		{
			Name:             "access-context",
			Builder:          NewTestBuilder().WithScope(scopeCtx).WithContext(context.WithValue(context.Background(), "foo", "bar")),
			ExpectedSQL:      "SELECT \"foos\".* FROM \"foos\" WHERE \"a\" = ?",
			ExpectedBindings: []any{"bar"},
		},
	})
}
