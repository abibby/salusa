package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/database/insert"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/models"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

func TestDI(t *testing.T) {
	type Model struct {
		models.BaseModel
		ID   int    `db:"id,primary"`
		Name string `db:"name"`
	}
	r := dbtest.NewRunner(func() (*sqlx.DB, error) {
		db, err := sqlx.Open("sqlite3", ":memory:")
		if err != nil {
			return nil, err
		}
		c, err := migrate.CreateFromModel(&Model{})
		if err != nil {
			return nil, err
		}
		err = c.Run(context.Background(), db)
		if err != nil {
			return nil, err
		}
		return db, nil
	})
	r.Run(t, "fetches model", func(t *testing.T, tx *sqlx.Tx) {

		err := insert.Save(tx, &Model{ID: 1, Name: "test"})
		assert.NoError(t, err)

		type Req struct {
			Model *Model `di:"model"`
		}
		req := &Req{}
		r := httptest.NewRequest("GET", "/model/1", http.NoBody)
		r = mux.SetURLVars(r, map[string]string{"model": "1"})
		r = SetTestTx(r, tx)
		err = di(req, r)

		assert.NoError(t, err)
		if assert.NotNil(t, req.Model) {
			assert.Equal(t, "test", req.Model.Name)
		}
	})

	r.Run(t, "doesn't fetch wrong model", func(t *testing.T, tx *sqlx.Tx) {

		err := insert.Save(tx, &Model{ID: 1, Name: "test"})
		assert.NoError(t, err)

		type Req struct {
			Model *Model `di:"model"`
		}
		req := &Req{}
		r := httptest.NewRequest("GET", "/model/1", http.NoBody)
		r = mux.SetURLVars(r, map[string]string{"model": "2"})
		r = SetTestTx(r, tx)
		err = di(req, r)

		assert.NoError(t, err)
		assert.Nil(t, req.Model)
	})

	r.Run(t, "missing transaction", func(t *testing.T, tx *sqlx.Tx) {
		type Req struct {
			Model *Model `di:"model"`
		}
		req := &Req{}
		r := httptest.NewRequest("GET", "/model/1", http.NoBody)
		r = mux.SetURLVars(r, map[string]string{"model": "2"})
		err := di(req, r)

		assert.ErrorIs(t, err, ErrNoTx)
	})

	r.Run(t, "fetche non model", func(t *testing.T, tx *sqlx.Tx) {
		type Req struct {
			Model string `di:"model"`
		}
		req := &Req{}
		r := httptest.NewRequest("GET", "/model/1", http.NoBody)
		r = mux.SetURLVars(r, map[string]string{"model": "1"})
		r = SetTestTx(r, tx)
		err := di(req, r)

		assert.ErrorIs(t, err, ErrNotModel)
	})
}
