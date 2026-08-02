package migrate_test

import (
	"testing"

	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMustMigrateModel(t *testing.T) {
	test.RunNoTx(t, "success", func(t *testing.T, db *sqlx.DB) {
		type MustMigrateModelModel struct {
			model.BaseModel
			ID int `db:"id,primary"`
		}
		assert.NotPanics(t, func() {
			migrate.MustMigrateModel(db, &MustMigrateModelModel{})
		})
	})

	test.RunNoTx(t, "panics on invalid field", func(t *testing.T, db *sqlx.DB) {
		type MustMigrateModelErrorModel struct {
			model.BaseModel
			ID    int `db:"id,primary"`
			Value complex128
		}
		assert.Panics(t, func() {
			migrate.MustMigrateModel(db, &MustMigrateModelErrorModel{})
		})
	})
}
