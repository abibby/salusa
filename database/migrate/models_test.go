package migrate_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRunModelCreate(t *testing.T) {
	test.Run(t, "success", func(t *testing.T, db *sqlx.Tx) {
		type RunModelCreateModel struct {
			model.BaseModel
			ID int `db:"id,primary"`
		}
		err := migrate.RunModelCreate(context.Background(), db, &RunModelCreateModel{})
		assert.NoError(t, err)
	})

	test.Run(t, "error", func(t *testing.T, db *sqlx.Tx) {
		type RunModelCreateErrorModel struct {
			model.BaseModel
			ID    int `db:"id,primary"`
			Value complex128
		}
		err := migrate.RunModelCreate(context.Background(), db, &RunModelCreateErrorModel{})
		assert.Error(t, err)
	})
}
