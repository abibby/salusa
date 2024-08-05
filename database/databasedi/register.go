package databasedi

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusaconfig"
	"github.com/jmoiron/sqlx"
)

type dbDeps struct {
	Cfg salusaconfig.Config `inject:""`
	Log *slog.Logger        `inject:""`
}

func RegisterFromConfig(migrations *migrate.Migrations) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		di.RegisterLazySingletonWith(ctx, func(deps *dbDeps) (*sqlx.DB, error) {
			var cfgAny any = deps.Cfg
			cfger, ok := cfgAny.(database.DBConfiger)
			if !ok {
				return nil, fmt.Errorf("config not instance of dialects.DBConfiger")
			}
			dbcfg := cfger.DBConfig()

			dbcfg.SetDialect()
			db, err := sqlx.Open(dbcfg.DriverName(), dbcfg.DataSourceName())
			if err != nil {
				return nil, fmt.Errorf("databasedi.RegisterFromConfig: open database: %w", err)
			}

			if migrations != nil {
				err = migrations.Up(ctx, db)
				if err != nil {
					return nil, fmt.Errorf("databasedi.RegisterFromConfig: migrate database: %w", err)
				}
			}
			deps.Log.Info("database ready")

			return db, nil
		})

		err := RegisterTransactions(nil)(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

func Register(db *sqlx.DB) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		di.RegisterSingleton(ctx, func() *sqlx.DB {
			return db
		})
		err := RegisterTransactions(nil)(ctx)
		if err != nil {
			return err
		}
		return nil
	}
}

func RegisterTransactions(mtx sync.Locker) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		di.RegisterWith(ctx, func(ctx context.Context, tag string, db *sqlx.DB) (database.Read, error) {
			return database.NewRead(ctx, mtx, db), nil
		})
		di.RegisterWith(ctx, func(ctx context.Context, tag string, db *sqlx.DB) (database.Update, error) {
			return database.NewUpdate(ctx, mtx, db), nil
		})
		return nil
	}
}
