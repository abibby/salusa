package dbtest_test

import (
	"testing"

	"abibby.com/salusa/database/dbtest"
	"abibby.com/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func newRunner() *dbtest.Runner {
	return dbtest.NewRunner(func() (*sqlx.DB, error) {
		return sqlx.Open("sqlite3", ":memory:")
	})
}

func TestUpdateRead(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		u := dbtest.Update(tx)
		err := u(func(tx *sqlx.Tx) error { return nil })
		assert.NoError(t, err)

		r := dbtest.Read(tx)
		err = r(func(tx *sqlx.Tx) error { return nil })
		assert.NoError(t, err)
	})
}
