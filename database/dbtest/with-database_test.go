package dbtest_test

import (
	_ "github.com/mattn/go-sqlite3"
)

// func newSQLiteRunner() *dbtest.Runner {
// 	return dbtest.NewRunner(func() (*sqlx.DB, error) {
// 		return sqlx.Open("sqlite3", ":memory:")
// 	})
// }

// // func TestRunNoTx(t *testing.T) {
// // 	runner := newSQLiteRunner()

// // 	ran := false
// // 	ok := runner.RunNoTx(t, "no tx", func(t *testing.T, db *sqlx.DB) {
// // 		ran = true
// // 		assert.NotNil(t, db)
// // 	})
// // 	assert.True(t, ok)
// // 	assert.True(t, ran)
// // }

// func TestRunNoTxOpenError(t *testing.T) {
// 	runner := dbtest.NewRunner(func() (*sqlx.DB, error) {
// 		return sqlx.Open("bogus-driver", ":memory:")
// 	})

// 	fakeT := &testing.T{}
// 	ok := runner.RunNoTx(fakeT, "open error", func(t *testing.T, db *sqlx.DB) {
// 		t.Fatal("should not run")
// 	})
// 	assert.False(t, ok)
// }

// func TestRunBenchmarkAndRunBenchmarkNoTx(t *testing.T) {
// 	runner := newSQLiteRunner()

// 	ranTx := false
// 	result := testing.Benchmark(func(b *testing.B) {
// 		runner.RunBenchmark(b, "bench tx", func(b *testing.B, tx *sqlx.Tx) {
// 			ranTx = true
// 			assert.NotNil(b, tx)
// 		})
// 	})
// 	assert.True(t, ranTx)
// 	assert.NotNil(t, result)

// 	ranNoTx := false
// 	result = testing.Benchmark(func(b *testing.B) {
// 		runner.RunBenchmarkNoTx(b, "bench no tx", func(b *testing.B, db *sqlx.DB) {
// 			ranNoTx = true
// 			assert.NotNil(b, db)
// 		})
// 	})
// 	assert.True(t, ranNoTx)
// 	assert.NotNil(t, result)
// }
