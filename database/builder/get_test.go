package builder_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		const insert = "INSERT INTO foos (id, name) values (?,?)"
		_, err := tx.ExecContext(context.Background(), insert, 1, "test1")
		assert.NoError(t, err)
		_, err = tx.ExecContext(context.Background(), insert, 2, "test2")
		assert.NoError(t, err)

		foos, err := builder.From[*test.Foo]().Get(tx)
		assert.NoError(t, err)
		assertJsonEqual(t, `[
			{"id":1,"name":"test1","bar":null,"bars":null},
			{"id":2,"name":"test2","bar":null,"bars":null}
		]`, foos)
	})
}

func TestFirst(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		const insert = "INSERT INTO foos (id, name) values (?,?)"
		_, err := tx.ExecContext(context.Background(), insert, 1, "test1")
		assert.NoError(t, err)
		_, err = tx.ExecContext(context.Background(), insert, 2, "test2")
		assert.NoError(t, err)

		foo, err := builder.From[*test.Foo]().First(tx)
		assert.NoError(t, err)
		assertJsonEqual(t, `{
			"id":1,
			"name":"test1",
			"bar":null,
			"bars":null
		}`, foo)
	})
}

func TestGet_with_scope_and_context(t *testing.T) {
	scopeCtx := &builder.Scope{
		Name: "ctx",
		Apply: func(b *builder.SubBuilder) *builder.SubBuilder {
			return b.Where("id", "=", b.Context().Value("foo"))
		},
	}
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.WithValue(context.Background(), "foo", 2)

		MustSave(tx, &test.Foo{ID: 1, Name: "foo1"})
		MustSave(tx, &test.Foo{ID: 2, Name: "foo2"})

		foos, err := NewTestBuilder().
			WithScope(scopeCtx).
			Where("name", "like", "foo%").
			WithContext(ctx).
			Get(tx)
		assert.NoError(t, err)
		assertJsonEqual(t, `[{
			"id":2,
			"name":"foo2",
			"bar":null,
			"bars":null
		}]`, foos)
	})

}

func TestGet_returns_empty_array_with_no_results(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		foos, err := NewTestBuilder().Get(tx)
		assert.NoError(t, err)
		assertJsonEqual(t, `[]`, foos)
	})
}

func TestFirst_returns_nil_with_no_results(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		foo, err := NewTestBuilder().First(tx)
		assert.NoError(t, err)
		assert.Nil(t, foo)
	})
}

func TestEach(t *testing.T) {
	test.Run(t, "runs for every row", func(t *testing.T, tx *sqlx.Tx) {
		MustSave(tx, &test.Foo{ID: 1, Name: "foo1"})
		MustSave(tx, &test.Foo{ID: 2, Name: "foo2"})

		i := 0
		err := NewTestBuilder().Each(tx, func(v *test.Foo) error {
			i++
			assert.Equal(t, i, v.ID)
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 2, i)
	})
}
