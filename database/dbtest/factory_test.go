package dbtest_test

import (
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestFactory(t *testing.T) {
	fooFactory := dbtest.NewFactory(func() *test.Foo {
		return &test.Foo{
			Name: "foo",
		}
	})
	test.Run(t, "create", func(t *testing.T, tx *sqlx.Tx) {
		f := fooFactory.Create(tx)
		assert.Equal(t, "foo", f.Name)

		dbF, err := builder.From[*test.Foo]().Find(tx, f.ID)
		assert.NoError(t, err)
		assert.Equal(t, f, dbF)
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
			State(func(f *test.Foo) *test.Foo {
				f.Name = "bar"
				return f
			}).
			Create(tx)
		assert.Equal(t, "bar", f.Name)
	})
}
