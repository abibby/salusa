package auth_test

import (
	"testing"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/di"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Run("register autoincrement", func(t *testing.T) {
		db := sqlx.MustOpen("sqlite3", ":memory:")
		defer db.Close()
		migrate.MustMigrateModel(db, &AutoIncrementUser{})
		createdUser := &AutoIncrementUser{
			Username:     "user",
			PasswordHash: []byte{},
		}
		err := model.Save(db, createdUser)
		assert.NoError(t, err)

		ctx := di.TestDependencyProviderContext()
		_ = databasedi.Register(db)(ctx)
		err = auth.Register[*AutoIncrementUser](ctx)
		assert.NoError(t, err)

		ctx = auth.SetClaims(ctx, &auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: createdUser.GetID(),
			},
		})

		u, err := di.Resolve[*AutoIncrementUser](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, u)
	})

	t.Run("register uuid", func(t *testing.T) {
		db := sqlx.MustOpen("sqlite3", ":memory:")
		defer db.Close()
		migrate.MustMigrateModel(db, &auth.UsernameUser{})
		createdUser := &auth.UsernameUser{
			ID:           uuid.New(),
			Username:     "user",
			PasswordHash: []byte{},
		}
		err := model.Save(db, createdUser)
		assert.NoError(t, err)

		ctx := di.TestDependencyProviderContext()
		_ = databasedi.Register(db)(ctx)
		err = auth.Register[*auth.UsernameUser](ctx)
		assert.NoError(t, err)

		ctx = auth.SetClaims(ctx, &auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: createdUser.GetID(),
			},
		})

		u, err := di.Resolve[*auth.UsernameUser](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, u)
	})

	t.Run("no claims", func(t *testing.T) {
		db := sqlx.MustOpen("sqlite3", ":memory:")
		defer db.Close()
		migrate.MustMigrateModel(db, &auth.UsernameUser{})
		createdUser := &auth.UsernameUser{
			ID:           uuid.New(),
			Username:     "user",
			PasswordHash: []byte{},
		}
		err := model.Save(db, createdUser)
		assert.NoError(t, err)

		ctx := di.TestDependencyProviderContext()
		_ = databasedi.Register(db)(ctx)
		err = auth.Register[*auth.UsernameUser](ctx)
		assert.NoError(t, err)

		u, err := di.Resolve[*auth.UsernameUser](ctx)
		assert.ErrorIs(t, err, auth.Err401Unauthorized)
		assert.Nil(t, u)
	})
}
