package databasedi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	db := sqlx.MustOpen("sqlite3", ":memory:")
	databasedi.Register(db)
	m := databasedi.Middleware()

	runs := 0

	handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := di.Resolve[*sqlx.Tx](r.Context())
		assert.NoError(t, err)
		assert.NotNil(t, tx)
		runs++
	}))

	handler.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "http://example.com", nil),
	)

	assert.Equal(t, 1, runs)

	t.Run("tx commits", func(t *testing.T) {
		t.Fail()
	})

	t.Run("tx rolls back", func(t *testing.T) {
		t.Fail()
	})
}
