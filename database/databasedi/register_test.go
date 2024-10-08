package databasedi_test

import (
	"testing"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	cfg := sqlite.NewConfig(":memory:")
	t.Run("db", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		db := sqlx.MustOpen(cfg.DriverName(), cfg.DataSourceName())
		defer db.Close()
		_ = databasedi.Register(db)(ctx)

		newDB, err := di.Resolve[*sqlx.DB](ctx)
		assert.NoError(t, err)
		assert.Same(t, db, newDB)
	})

	t.Run("tx read", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		db := sqlx.MustOpen(cfg.DriverName(), cfg.DataSourceName())
		defer db.Close()
		_ = databasedi.Register(db)(ctx)

		read, err := di.Resolve[database.Read](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, read)
	})

	t.Run("tx update", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		db := sqlx.MustOpen(cfg.DriverName(), cfg.DataSourceName())
		defer db.Close()
		_ = databasedi.Register(db)(ctx)

		update, err := di.Resolve[database.Update](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, update)
	})

	t.Run("tx", func(t *testing.T) {
		ctx := di.TestDependencyProviderContext()
		db := sqlx.MustOpen(cfg.DriverName(), cfg.DataSourceName())
		defer db.Close()
		_ = databasedi.Register(db)(ctx)

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
