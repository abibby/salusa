package builder_test

import (
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMulti(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		foos := []*test.Foo{
			{ID: 1},
			// {ID: 2},
			// {ID: 3},
		}
		for _, f := range foos {
			MustSave(tx, f)
			MustSave(tx, &test.Bar{ID: f.ID + 3, FooID: f.ID})
		}

		err := builder.Load(tx, foos, "Bar")
		assert.NoError(t, err)

		err = builder.Load(tx, foos, "Bars")
		assert.NoError(t, err)

		for _, foo := range foos {
			assert.True(t, foo.Bar.Loaded(), "bar not loaded")

			assert.True(t, foo.Bars.Loaded(), "bars not loaded")
		}
	})
}
