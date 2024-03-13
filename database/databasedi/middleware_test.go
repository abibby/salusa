package databasedi_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	dp := di.NewDependencyProvider()
	db := sqlx.MustOpen("sqlite3", ":memory:")
	databasedi.Register(dp, db)
	m := databasedi.Middleware()

	t.Run("injects database", func(t *testing.T) {

		runs := 0

		handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			injectedDB, err := di.Resolve[*sqlx.DB](r.Context(), dp)
			assert.NoError(t, err)
			assert.Same(t, db, injectedDB)
			runs++
		}))

		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com", nil),
		)

		assert.Equal(t, 1, runs)
	})

	t.Run("injects transaction", func(t *testing.T) {

		runs := 0

		handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := di.Resolve[*sqlx.Tx](r.Context(), dp)
			assert.NoError(t, err)
			assert.NotNil(t, tx)
			runs++
		}))

		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com", nil),
		)

		assert.Equal(t, 1, runs)
	})

	t.Run("tx commits", func(t *testing.T) {
		var tx *sqlx.Tx
		var err error
		handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err = di.Resolve[*sqlx.Tx](r.Context(), dp)
			assert.NoError(t, err)
			_, err := tx.Exec("create table tx_commits (id integer not null);")
			assert.NoError(t, err)
		}))

		ctx, done := context.WithCancel(context.Background())
		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com", nil).WithContext(ctx),
		)

		done()

		time.Sleep(time.Millisecond * 100)

		err = tx.Rollback()
		assert.ErrorIs(t, err, sql.ErrTxDone)

		tables := []string{}
		err = db.Select(&tables, `SELECT name FROM sqlite_schema WHERE type='table' AND name=?`, "tx_commits")
		assert.NoError(t, err)
		assert.Len(t, tables, 1)
	})

	t.Run("tx rolls back", func(t *testing.T) {
		var tx *sqlx.Tx
		var err error
		handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err = di.Resolve[*sqlx.Tx](r.Context(), dp)
			assert.NoError(t, err)
			_, err := tx.Exec("create table tx_rolls_back (id integer not null);")
			assert.NoError(t, err)
			w.WriteHeader(500)
		}))

		ctx, done := context.WithCancel(context.Background())
		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest("GET", "http://example.com", nil).WithContext(ctx),
		)

		done()

		time.Sleep(time.Millisecond * 100)

		err = tx.Rollback()
		assert.ErrorIs(t, err, sql.ErrTxDone)

		tables := []string{}
		err = db.Select(&tables, `SELECT name FROM sqlite_schema WHERE type='table' AND name=?`, "tx_rolls_back")
		assert.NoError(t, err)
		assert.Len(t, tables, 0)
	})

	t.Run("tx rolls back panic", func(t *testing.T) {
		var tx *sqlx.Tx
		var err error

		handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err = di.Resolve[*sqlx.Tx](r.Context(), dp)
			assert.NoError(t, err)
			_, err := tx.Exec("create table tx_rolls_back_panic (id integer not null);")
			assert.NoError(t, err)
			panic("error")
		}))

		ctx, done := context.WithCancel(context.Background())
		assert.Panics(t, func() {
			handler.ServeHTTP(
				httptest.NewRecorder(),
				httptest.NewRequest("GET", "http://example.com", nil).WithContext(ctx),
			)
		})

		done()

		time.Sleep(time.Millisecond * 100)

		err = tx.Rollback()
		assert.ErrorIs(t, err, sql.ErrTxDone)

		tables := []string{}
		err = db.Select(&tables, `SELECT name FROM sqlite_schema WHERE type='table' AND name=?`, "tx_rolls_back_panic")
		assert.NoError(t, err)
		assert.Len(t, tables, 0)
	})
}
