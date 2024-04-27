package auth_test

import (
	"bytes"
	"context"
	"embed"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/email/emailtest"
	"github.com/abibby/salusa/router/routertest"
	"github.com/abibby/salusa/view"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

type DevNull struct{}

func (d *DevNull) Write(p []byte) (n int, err error) {
	return len(p), nil
}

//go:embed emails/*
var emails embed.FS

var nullLogger = slog.New(slog.NewTextHandler(&DevNull{}, nil))
var usernameRoutes = auth.Routes(auth.NewUsernameUser, "reset-password")
var emailRoutes = auth.Routes(auth.NewEmailVerifiedUser, "reset-password")
var emailTemplates = view.NewViewTemplate(emails)

func TestAuthRoutesUserCreate(t *testing.T) {
	Run(t, "can create user", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		resp, err := usernameRoutes.UserCreate.Run(&auth.UserCreateRequest{
			Password: "pass",
			Update:   dbtest.Update(tx),
			Ctx:      ctx,
			Logger:   nullLogger,
			Request:  httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(`{"username":"user"}`)),
		})
		assert.NoError(t, err)
		assert.Equal(t, "user", resp.User.Username)

		u, err := builder.From[*auth.UsernameUser]().WithContext(ctx).Find(tx, resp.User.ID)
		assert.NoError(t, err)
		assert.Equal(t, u, resp.User)
		assert.NotNil(t, u.PasswordHash)
		// assert.False(t, u.Validated)
		// assert.NotZero(t, u.ValidationCode)
	})

	Run(t, "force lowercase usernames", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		resp, err := usernameRoutes.UserCreate.Run(&auth.UserCreateRequest{
			Password: "pass",
			Update:   dbtest.Update(tx),
			Ctx:      ctx,
			Logger:   nullLogger,
			Request:  httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(`{"username":"USER"}`)),
		})
		assert.NoError(t, err)
		assert.Equal(t, "user", resp.User.Username)

		u, err := builder.From[*auth.UsernameUser]().WithContext(ctx).Find(tx, resp.User.ID)
		assert.NoError(t, err)
		assert.Equal(t, u, resp.User)
		assert.NotNil(t, u.PasswordHash)
		// assert.False(t, u.Validated)
		// assert.NotZero(t, u.ValidationCode)
	})

	Run(t, "email verification", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		urlResolver := routertest.NewTestResolver()
		m := emailtest.NewTestMailer()
		resp, err := emailRoutes.UserCreate.Run(&auth.UserCreateRequest{
			Password: "pass",
			Update:   dbtest.Update(tx),
			Ctx:      ctx,
			Mailer:   m,
			Logger:   nullLogger,
			Request:  httptest.NewRequest("POST", "/user/create", bytes.NewBufferString(`{"username":"user@example.com"}`)),
			URL:      urlResolver,
			Template: emailTemplates,
		})
		assert.NoError(t, err)
		assert.Equal(t, "user@example.com", resp.User.Email)

		time.Sleep(time.Millisecond * 20)

		sent := m.EmailsSent()
		assert.Len(t, sent, 1)
		assert.Equal(t, []string{"user@example.com"}, sent[0].To)
		assert.Equal(t, "Verify your email", sent[0].Subject)
		assert.Contains(t, string(sent[0].HTMLBody), urlResolver.Resolve("auth.email.verify", "token", resp.User.LookupToken))

		u, err := builder.From[*auth.EmailVerifiedUser]().WithContext(ctx).Find(tx, resp.User.ID)
		assert.NoError(t, err)
		assert.Equal(t, u, resp.User)

		assert.False(t, u.Verified)
		assert.NotZero(t, u.LookupToken)
	})
}

func TestAuthRoutesLogin(t *testing.T) {

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
		u := &auth.UsernameUser{
			ID:           id,
			Username:     "user",
			PasswordHash: passwordHash,
		}
		err := model.Save(tx, u)
		assert.NoError(t, err)

		resp, err := usernameRoutes.Login.Run(&auth.LoginRequest{
			Username: "user",
			Password: password,
			Read:     dbtest.Read(tx),
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
		u := &auth.UsernameUser{
			ID:           uuid.New(),
			Username:     "user",
			PasswordHash: passwordHash,
		}
		err := model.Save(tx, u)
		assert.NoError(t, err)

		_, err = usernameRoutes.Login.Run(&auth.LoginRequest{
			Username: "user",
			Password: password,
			Read:     dbtest.Read(tx),
			Ctx:      ctx,
		})
		assert.ErrorIs(t, err, auth.ErrInvalidUserPass)
	})
	Run(t, "wrong user", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		u := &auth.UsernameUser{
			ID:           id,
			Username:     "user",
			PasswordHash: passwordHash,
		}
		err := model.Save(tx, u)
		assert.NoError(t, err)

		_, err = usernameRoutes.Login.Run(&auth.LoginRequest{
			Username: "not user",
			Password: password,
			Read:     dbtest.Read(tx),
			Ctx:      ctx,
		})
		assert.ErrorIs(t, err, auth.ErrInvalidUserPass)
	})
	Run(t, "wrong pass", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		u := &auth.UsernameUser{
			ID:           id,
			Username:     "user",
			PasswordHash: passwordHash,
		}
		err := model.Save(tx, u)
		assert.NoError(t, err)

		_, err = usernameRoutes.Login.Run(&auth.LoginRequest{
			Username: "user",
			Password: "not pass",
			Read:     dbtest.Read(tx),
			Ctx:      ctx,
		})
		assert.ErrorIs(t, err, auth.ErrInvalidUserPass)
	})
}

func TestAuthRoutesVerifyEmail(t *testing.T) {
	Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		urlResolver := routertest.NewTestResolver()
		token := "test"
		id := uuid.New()
		err := model.Save(tx, &auth.EmailVerifiedUser{
			ID:           id,
			Email:        "",
			PasswordHash: []byte{},
			LookupToken:  token,
			Verified:     false,
		})
		assert.NoError(t, err)

		resp, err := emailRoutes.VerifyEmail.Run(&auth.VerifyEmailRequest{
			Token:  token,
			Ctx:    ctx,
			Update: dbtest.Update(tx),
			URL:    urlResolver,
		})
		assert.NoError(t, err)

		u, err := builder.From[*auth.EmailVerifiedUser]().Find(tx, id)
		assert.NoError(t, err)
		assert.True(t, u.Verified)

		w := httptest.NewRecorder()
		resp.ServeHTTP(w, httptest.NewRequest("GET", "/", http.NoBody))
		assert.Equal(t, http.StatusFound, w.Result().StatusCode)
		assert.Equal(t, urlResolver.Resolve("login"), w.Result().Header.Get("Location"))
	})
}

func TestAuthRoutesResetPassword(t *testing.T) {
	Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		urlResolver := routertest.NewTestResolver()
		token := "lookup token"
		id := uuid.New()
		oldPasswordHash := []byte("old hash")
		err := model.Save(tx, &auth.EmailVerifiedUser{
			ID:           id,
			Email:        "",
			PasswordHash: oldPasswordHash,
			LookupToken:  token,
			Verified:     true,
		})
		assert.NoError(t, err)

		resp, err := emailRoutes.ResetPassword.Run(&auth.ResetPasswordRequest{
			Token:    token,
			Password: "new password",
			Ctx:      ctx,
			Update:   dbtest.Update(tx),
			URL:      urlResolver,
		})
		assert.NoError(t, err)

		u, err := builder.From[*auth.EmailVerifiedUser]().Find(tx, id)
		assert.NoError(t, err)
		assert.NotEqual(t, oldPasswordHash, u.PasswordHash)
		assert.Equal(t, u, resp.User)
	})
}

func TestAuthRoutesChangePassword(t *testing.T) {
	// Hashed password salted with the id
	id := uuid.MustParse("cae3c6b1-7ff1-4f23-9489-a9f6e82478f9")
	oldPassword := "pass"
	oldPasswordHash := []byte{
		0x24, 0x32, 0x61, 0x24, 0x30, 0x34, 0x24, 0x78, 0x4d, 0x65,
		0x30, 0x54, 0x66, 0x77, 0x4c, 0x75, 0x48, 0x79, 0x35, 0x78,
		0x64, 0x51, 0x76, 0x58, 0x6b, 0x59, 0x73, 0x4b, 0x2e, 0x36,
		0x34, 0x31, 0x70, 0x6c, 0x63, 0x6c, 0x69, 0x54, 0x43, 0x5a,
		0x51, 0x51, 0x55, 0x49, 0x71, 0x41, 0x72, 0x65, 0x77, 0x51,
		0x45, 0x4c, 0x6b, 0x43, 0x76, 0x6d, 0x6a, 0x62, 0x4d, 0x75,
	}

	Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		createdUser := &auth.UsernameUser{
			ID:           id,
			PasswordHash: oldPasswordHash,
		}
		err := model.Save(tx, createdUser)
		assert.NoError(t, err)

		resp, err := usernameRoutes.ChangePassword.Run(&auth.ChangePasswordRequest[*auth.UsernameUser]{
			OldPassword: oldPassword,
			NewPassword: "new password",
			User:        createdUser,
			Ctx:         ctx,
			Update:      dbtest.Update(tx),
		})
		assert.NoError(t, err)

		u, err := builder.From[*auth.UsernameUser]().Find(tx, id)
		assert.NoError(t, err)
		assert.NotEqual(t, oldPasswordHash, u.PasswordHash)
		assert.Equal(t, u, resp.User)
	})
}

func TestAuthRoutesRefresh(t *testing.T) {
	Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		createdUser := &auth.UsernameUser{
			ID:           uuid.New(),
			PasswordHash: []byte(""),
		}
		err := model.Save(tx, createdUser)
		assert.NoError(t, err)

		token, err := auth.GenerateTokenFrom(&auth.Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: createdUser.GetID(),
			},
			Type: auth.TypeRefresh,
		})
		assert.NoError(t, err)

		resp, err := usernameRoutes.Refresh.Run(&auth.RefreshRequest[*auth.UsernameUser]{
			RefreshToken: token,
			Ctx:          ctx,
			Read:         dbtest.Read(tx),
		})
		assert.NoError(t, err)
		assert.Equal(t, token, resp.RefreshToken)
		assert.Equal(t, "Bearer", resp.TokenType)
		assert.Equal(t, token, resp.RefreshToken)
		assert.Equal(t, 3600, resp.ExpiresIn)

		claims, err := auth.Parse(resp.AccessToken)
		assert.NoError(t, err)
		assert.Equal(t, createdUser.GetID(), claims.Subject)
		assert.Equal(t, auth.TypeAccess, claims.Type)
	})
}

func TestAuthRoutesForgotPassword(t *testing.T) {
	Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		ctx := context.Background()
		urlResolver := routertest.NewTestResolver()
		m := emailtest.NewTestMailer()

		id := uuid.New()
		err := model.Save(tx, &auth.EmailVerifiedUser{
			ID:           id,
			Email:        "user@example.com",
			PasswordHash: []byte{},
			Verified:     true,
		})
		assert.NoError(t, err)

		_, err = emailRoutes.ForgotPassword.Run(&auth.ForgotPasswordRequest{
			Email:    "user@example.com",
			Update:   dbtest.Update(tx),
			Ctx:      ctx,
			Mailer:   m,
			Logger:   nullLogger,
			URL:      urlResolver,
			Template: emailTemplates,
		})
		assert.NoError(t, err)

		time.Sleep(time.Millisecond * 20)

		u, err := builder.From[*auth.EmailVerifiedUser]().WithContext(ctx).Find(tx, id)
		assert.NoError(t, err)

		assert.NotZero(t, u.LookupToken)

		sent := m.EmailsSent()
		assert.Len(t, sent, 1)
		assert.Equal(t, []string{"user@example.com"}, sent[0].To)
		assert.Equal(t, "Password reset", sent[0].Subject)
		assert.Contains(t, string(sent[0].HTMLBody), urlResolver.Resolve("reset-password", "token", u.LookupToken))

	})
}
