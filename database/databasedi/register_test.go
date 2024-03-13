package databasedi_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Run("db", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		db := sqlx.MustOpen("sqlite3", ":memory:")
		databasedi.Register(dp, db)

		newDB, err := di.Resolve[*sqlx.DB](context.Background(), dp)
		assert.NoError(t, err)
		assert.Same(t, db, newDB)
	})

	t.Run("tx", func(t *testing.T) {
		dp := di.NewDependencyProvider()
		db := sqlx.MustOpen("sqlite3", ":memory:")
		databasedi.Register(dp, db)

		tx, err := di.Resolve[*sqlx.Tx](context.Background(), dp)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})
}
