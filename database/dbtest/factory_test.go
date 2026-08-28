package dbtest_test

import (
	"testing"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/dbtest"
	"abibby.com/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	fooFactory := dbtest.NewFactory(func(tx database.DB) *test.Foo {
		return &test.Foo{
			Name: "foo",
		}
	})
	test.Run(t, "create", func(t *testing.T, tx *sqlx.Tx) {
		f := fooFactory.Create(tx)
		assert.Equal(t, "foo", f.Name)

		dbF, err := builder.From[*test.Foo]().Find(tx, f.ID)
		if assert.NoError(t, err) {
			assert.Equal(t, f, dbF)
		}
	})
	test.Run(t, "count", func(t *testing.T, tx *sqlx.Tx) {
		foos := fooFactory.Count(4).Create(tx)
		assert.Len(t, foos, 4)
		for _, f := range foos {
			assert.Equal(t, "foo", f.Name)
		}
	})
	test.Run(t, "state", func(t *testing.T, tx *sqlx.Tx) {
		f := fooFactory.
			State(func(f *test.Foo) {
				f.Name = "bar"
			}).
			Create(tx)
		assert.Equal(t, "bar", f.Name)
	})
}
