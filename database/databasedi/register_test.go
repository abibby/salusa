package databasedi_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusaconfig"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	cfg database.Config
}

func (c *testConfig) DBConfig() database.Config { return c.cfg }
func (c *testConfig) GetHTTPPort() int          { return 8080 }
func (c *testConfig) GetBaseURL() string        { return "http://localhost" }

type plainConfig struct{}

func (plainConfig) GetHTTPPort() int   { return 8080 }
func (plainConfig) GetBaseURL() string { return "http://localhost" }

type badConfig struct{}

func (badConfig) DriverName() string     { return "no-such-driver" }
func (badConfig) DataSourceName() string { return "" }
func (badConfig) SetDialect()            {}

type badConfiger struct{}

func (badConfiger) DBConfig() database.Config { return badConfig{} }
func (badConfiger) GetHTTPPort() int          { return 8080 }
func (badConfiger) GetBaseURL() string        { return "http://localhost" }

func registerConfig(ctx context.Context, cfg salusaconfig.Config) {
	di.RegisterSingleton[salusaconfig.Config](ctx, func() salusaconfig.Config { return cfg })
	di.RegisterSingleton[*slog.Logger](ctx, func() *slog.Logger { return slog.Default() })
}

func TestRegisterFromConfig(t *testing.T) {
	ctx := di.TestDependencyProviderContext()
	registerConfig(ctx, &testConfig{cfg: sqlite.NewConfig(":memory:")})

	err := databasedi.RegisterFromConfig(migrate.New())(ctx)
	assert.NoError(t, err)

	db, err := di.Resolve[*sqlx.DB](ctx)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	db.Close()
}

func TestRegisterFromConfig_not_db_configer(t *testing.T) {
	ctx := di.TestDependencyProviderContext()
	registerConfig(ctx, plainConfig{})

	err := databasedi.RegisterFromConfig(nil)(ctx)
	assert.NoError(t, err)

	_, err = di.Resolve[*sqlx.DB](ctx)
	assert.Error(t, err)
}

func TestRegisterFromConfig_open_error(t *testing.T) {
	ctx := di.TestDependencyProviderContext()
	registerConfig(ctx, badConfiger{})

	err := databasedi.RegisterFromConfig(nil)(ctx)
	assert.NoError(t, err)

	_, err = di.Resolve[*sqlx.DB](ctx)
	assert.Error(t, err)
}

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
