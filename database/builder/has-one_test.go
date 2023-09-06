package builder_test

import (
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type HasOneFoo struct {
	test.Foo
	BadLocal   *builder.HasOne[*test.Bar] `db:"-" local:"bad_key"`
	BadForeign *builder.HasOne[*test.Bar] `db:"-" foreign:"bad_key"`
}

func TestHasOneLoad(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		foos := []*test.Foo{
			{ID: 1},
			{ID: 2},
			{ID: 3},
		}
		for _, f := range foos {
			MustSave(tx, f)
		}
		MustSave(tx, &test.Bar{ID: 4, FooID: 1})
		MustSave(tx, &test.Bar{ID: 5, FooID: 2})
		MustSave(tx, &test.Bar{ID: 6, FooID: 3})

		err := builder.Load(tx, foos, "Bar")
		assert.NoError(t, err)

		for _, foo := range foos {
			assert.True(t, foo.Bar.Loaded())
			bar, ok := foo.Bar.Value()
			assert.True(t, ok)
			assert.Equal(t, foo.ID+3, bar.ID)
			assert.Equal(t, foo.ID, bar.FooID)
		}
	})
}

func TestHasOne_json_marshal(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		f := &test.Foo{ID: 1}
		MustSave(tx, f)
		MustSave(tx, &test.Bar{ID: 4, FooID: 1})

		err := builder.Load(tx, f, "Bar")
		assert.NoError(t, err)

		assertJsonEqual(t, `{
			"id":1,
			"name":"",
			"bar":{"id":4,"foo_id":1,"foo":null},
			"bars":null
		}`, f)

	})

}

func TestHasOne_deep(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		f := &test.Foo{ID: 1}
		MustSave(tx, f)
		MustSave(tx, &test.Bar{ID: 4, FooID: 1})

		err := builder.Load(tx, f, "Bar.Foo")
		assert.NoError(t, err)

		assertJsonEqual(t, `{
			"id":1,
			"name":"",
			"bar":{
				"id":4,
				"foo_id":1,
				"foo": {
					"id":1,
					"name":"",
					"bar":null,
					"bars":null
				}
			},
			"bars":null
		}`, f)

	})

}

func TestHasOne_invalid_local_key(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		f := &HasOneFoo{Foo: test.Foo{ID: 1}}
		MustSave(tx, f)
		MustSave(tx, &test.Bar{ID: 4, FooID: 1})

		err := builder.Load(tx, f, "BadLocal")
		assert.ErrorIs(t, err, builder.ErrMissingField)
	})

}

func TestHasOne_invalid_foreign_key(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		f := &HasOneFoo{Foo: test.Foo{ID: 1}}
		MustSave(tx, f)
		MustSave(tx, &test.Bar{ID: 4, FooID: 1})

		err := builder.Load(tx, f, "BadForeign")
		assert.ErrorIs(t, err, builder.ErrMissingField)
	})

}

func BenchmarkHasOneLoad(b *testing.B) {
	test.RunBenchmark(b, "", func(t *testing.B, tx *sqlx.Tx) {
		foos := make([]*test.Foo, 100)
		for i := 0; i < 100; i++ {
			f := &test.Foo{ID: i}
			foos[i] = f
			MustSave(tx, f)
			MustSave(tx, &test.Bar{ID: i, FooID: i})
		}
		for i := 0; i < b.N; i++ {
			builder.Load(tx, foos, "Bars")
		}

	})
}
