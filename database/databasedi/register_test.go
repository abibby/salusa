package databasedi_test

import (
	"testing"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Run("db", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		db := sqlx.MustOpen("sqlite3", ":memory:")
		defer db.Close()
		databasedi.Register(ctx, db)

		newDB, err := di.Resolve[*sqlx.DB](ctx)
		assert.NoError(t, err)
		assert.Same(t, db, newDB)
	})

	t.Run("tx read", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		db := sqlx.MustOpen("sqlite3", ":memory:")
		defer db.Close()
		databasedi.Register(ctx, db)

		read, err := di.Resolve[database.Read](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, read)
	})

	t.Run("tx update", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		db := sqlx.MustOpen("sqlite3", ":memory:")
		defer db.Close()
		databasedi.Register(ctx, db)

		update, err := di.Resolve[database.Update](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, update)
	})

	t.Run("tx", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		db := sqlx.MustOpen("sqlite3", ":memory:")
		defer db.Close()
		databasedi.Register(ctx, db)

		update, err := di.Resolve[database.Update](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, update)

		run := 0
		err = update(func(tx *sqlx.Tx) error {
			run++
			assert.NotNil(t, tx)
			return nil
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, run)
	})
}
