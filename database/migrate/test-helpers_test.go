package migrate_test

import (
	"testing"

	"abibby.com/salusa/database/migrate"
	"abibby.com/salusa/database/model"
	"abibby.com/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestMustMigrateModel(t *testing.T) {
	test.Run(t, "success", func(t *testing.T, db *sqlx.Tx) {
		type MustMigrateModelModel struct {
			model.BaseModel
			ID int `db:"id,primary"`
		}
		assert.NotPanics(t, func() {
			migrate.MustMigrateModel(db, &MustMigrateModelModel{})
		})
	})

	test.Run(t, "panics on invalid field", func(t *testing.T, db *sqlx.Tx) {
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
