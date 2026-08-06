package dbtest_test

import (
	"testing"

	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/internal/test"
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

func TestRunNoTx(t *testing.T) {
	ran := false
	newRunner().RunNoTx(t, "", func(t *testing.T, db *sqlx.DB) {
		ran = true
		assert.NotNil(t, db)
	})
	assert.True(t, ran)
}

func TestRunBenchmarks(t *testing.T) {
	testing.Benchmark(func(b *testing.B) {
		ran := false
		newRunner().RunBenchmark(b, "", func(t *testing.B, tx *sqlx.Tx) {
			ran = true
			assert.NotNil(t, tx)
		})
		assert.True(t, ran)
	})
	testing.Benchmark(func(b *testing.B) {
		ran := false
		newRunner().RunBenchmarkNoTx(b, "", func(t *testing.B, db *sqlx.DB) {
			ran = true
			assert.NotNil(t, db)
		})
		assert.True(t, ran)
	})
}
