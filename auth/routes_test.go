package auth_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/email/emailtest"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type EmailVerifiedUser struct {
	auth.BaseUser
	auth.MustVerifyEmail
}

var runner = dbtest.NewRunner(func() (*sqlx.DB, error) {
	sqlite.UseSQLite()
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	err = migrate.RunModelCreate(ctx, db, &auth.BaseUser{}, &EmailVerifiedUser{})
	if err != nil {
		return nil, err
	}
	return db, nil
})

var Run = runner.Run

func TestUserCreate(t *testing.T) {
	routes := auth.Routes(func(base *auth.BaseUser) *auth.BaseUser {
		return base
	})
	Run(t, "can create user", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		resp, err := routes.UserCreate.Run(&auth.UserCreateRequest{
			Username: "user",
			Password: "pass",
			Tx:       tx,
			Ctx:      ctx,
		})
		assert.NoError(t, err)
		assert.Equal(t, "user", resp.User.Username)

		u, err := builder.From[*auth.BaseUser]().WithContext(ctx).Find(tx, resp.User.ID)
		assert.NoError(t, err)
		assert.Equal(t, u, resp.User)
		assert.NotNil(t, u.PasswordHash)
		// assert.False(t, u.Validated)
		// assert.NotZero(t, u.ValidationCode)
	})
}

func TestLogin(t *testing.T) {
	routes := auth.Routes(func(base *auth.BaseUser) *auth.BaseUser {
		return base
	})

	// Hashed password salted with the id
	id := uuid.MustParse("cae3c6b1-7ff1-4f23-9489-a9f6e82478f9")
	password := "pass"
	passwordHash := []byte{
		0x24, 0x32, 0x61, 0x24, 0x30, 0x34, 0x24, 0x78, 0x4d, 0x65,
		0x30, 0x54, 0x66, 0x77, 0x4c, 0x75, 0x48, 0x79, 0x35, 0x78,
		0x64, 0x51, 0x76, 0x58, 0x6b, 0x59, 0x73, 0x4b, 0x2e, 0x36,
		0x34, 0x31, 0x70, 0x6c, 0x63, 0x6c, 0x69, 0x54, 0x43, 0x5a,
		0x51, 0x51, 0x55, 0x49, 0x71, 0x41, 0x72, 0x65, 0x77, 0x51,
		0x45, 0x4c, 0x6b, 0x43, 0x76, 0x6d, 0x6a, 0x62, 0x4d, 0x75,
	}

	Run(t, "can login", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		u := &auth.BaseUser{
			ID:           id,
			Username:     "user",
			PasswordHash: passwordHash,
		}
		err := model.Save(tx, u)
		assert.NoError(t, err)

		resp, err := routes.Login.Run(&auth.LoginRequest{
			Username: "user",
			Password: password,
			Tx:       tx,
			Ctx:      ctx,
		})
		assert.NoError(t, err)
		assert.NotZero(t, resp.AccessToken)
		assert.NotZero(t, resp.RefreshToken)
		assert.Equal(t, "Bearer", resp.TokenType)
		assert.Equal(t, 3600, resp.ExpiresIn)
	})

	Run(t, "password is salted", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		u := &auth.BaseUser{
			ID:           uuid.New(),
			Username:     "user",
			PasswordHash: passwordHash,
		}
		err := model.Save(tx, u)
		assert.NoError(t, err)

		_, err = routes.Login.Run(&auth.LoginRequest{
			Username: "user",
			Password: password,
			Tx:       tx,
			Ctx:      ctx,
		})
		assert.ErrorIs(t, err, auth.ErrInvalidUserPass)
	})
	Run(t, "wrong user", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		u := &auth.BaseUser{
			ID:           id,
			Username:     "user",
			PasswordHash: passwordHash,
		}
		err := model.Save(tx, u)
		assert.NoError(t, err)

		_, err = routes.Login.Run(&auth.LoginRequest{
			Username: "not user",
			Password: password,
			Tx:       tx,
			Ctx:      ctx,
		})
		assert.ErrorIs(t, err, auth.ErrInvalidUserPass)
	})
	Run(t, "wrong pass", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		u := &auth.BaseUser{
			ID:           id,
			Username:     "user",
			PasswordHash: passwordHash,
		}
		err := model.Save(tx, u)
		assert.NoError(t, err)

		_, err = routes.Login.Run(&auth.LoginRequest{
			Username: "user",
			Password: "not pass",
			Tx:       tx,
			Ctx:      ctx,
		})
		assert.ErrorIs(t, err, auth.ErrInvalidUserPass)
	})
}

func TestEmailVerification(t *testing.T) {
	routes := auth.Routes(func(base *auth.BaseUser) *EmailVerifiedUser {
		return &EmailVerifiedUser{BaseUser: *base}
	})

	Run(t, "email sent", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		m := emailtest.NewTestMailer()
		resp, err := routes.UserCreate.Run(&auth.UserCreateRequest{
			Username: "user",
			Password: "pass",
			Tx:       tx,
			Ctx:      ctx,
			Mailer:   m,
		})
		assert.NoError(t, err)
		assert.Equal(t, "user", resp.User.Username)

		sent := m.EmailsSent()
		assert.Len(t, sent, 1)
		assert.Equal(t, "", sent[0].To)
		assert.Equal(t, "", sent[0].Body)
	})
}
