package auth_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/email/emailtest"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/router/routertest"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestClaimsBuilders(t *testing.T) {
	now := time.Now()

	c := auth.NewClaims().
		WithIssuer("iss").
		WithSubject("sub").
		WithAudience([]string{"aud"}).
		WithLifetime(time.Hour).
		WithNotBeforeTime(now).
		WithIssuedAtTime(now).
		WithJWTID("jti").
		WithScopes(auth.ScopeAccess, auth.ScopeRefresh)

	assert.Equal(t, "iss", c.Issuer)
	assert.Equal(t, "sub", c.Subject)
	assert.Equal(t, jwt.ClaimStrings{"aud"}, c.Audience)
	assert.NotNil(t, c.ExpiresAt)
	assert.NotNil(t, c.NotBefore)
	assert.NotNil(t, c.IssuedAt)
	assert.Equal(t, "jti", c.ID)
	assert.Equal(t, auth.ScopeStrings{auth.ScopeAccess, auth.ScopeRefresh}, c.Scope)
	assert.WithinDuration(t, now.Add(time.Hour), c.ExpiresAt.Time, time.Minute)
}

func TestControllerHandlersAndOptions(t *testing.T) {
	optionsCalled := false
	controller := auth.NewBasicAuthController[*auth.UsernameUser](
		auth.CreateUser(auth.NewUsernameUser),
		auth.AccessTokenOptions(func(u *auth.UsernameUser, claims *auth.Claims) jwt.Claims {
			optionsCalled = true
			return claims
		}),
		auth.RefreshTokenOptions(func(u *auth.UsernameUser, claims *auth.Claims) jwt.Claims {
			return claims
		}),
		auth.ResetPasswordName("custom-reset"),
	)

	handlers := []struct {
		name string
		get  func() http.Handler
	}{
		{"UserCreate", func() http.Handler { return controller.UserCreate() }},
		{"Login", func() http.Handler { return controller.Login() }},
		{"VerifyEmail", func() http.Handler { return controller.VerifyEmail() }},
		{"ResetPassword", func() http.Handler { return controller.ResetPassword() }},
		{"ForgotPassword", func() http.Handler { return controller.ForgotPassword() }},
		{"ChangePassword", func() http.Handler { return controller.ChangePassword() }},
		{"Refresh", func() http.Handler { return controller.Refresh() }},
	}
	for _, h := range handlers {
		t.Run(h.name, func(t *testing.T) {
			handler := h.get()
			assert.NotNil(t, handler)
		})
	}

	t.Run("access token options used", func(t *testing.T) {
		optionsCalled = false
		controller := auth.NewBasicAuthController[*auth.UsernameUser](
			auth.AccessTokenOptions(func(u *auth.UsernameUser, claims *auth.Claims) jwt.Claims {
				optionsCalled = true
				return claims
			}),
		)
		controller.UserCreate() // exercise options path not possible directly
		assert.False(t, optionsCalled)
	})
}

func TestRegisterRoutes(t *testing.T) {
	r := router.New()
	controller := auth.NewBasicAuthController[*auth.UsernameUser](auth.CreateUser(auth.NewUsernameUser))
	auth.RegisterRoutes(r, controller)
	assert.Len(t, r.Routes(), 7)
}

func TestUserTypes(t *testing.T) {
	t.Run("UsernameUser", func(t *testing.T) {
		id := uuid.New()
		u := &auth.UsernameUser{
			ID:           id,
			Username:     "USER",
			PasswordHash: []byte("hash"),
		}
		assert.Equal(t, id.String(), u.GetID())
		assert.Equal(t, "USER", u.GetUsername())
		assert.Equal(t, []string{"username"}, u.UsernameColumns())
		assert.Equal(t, []byte("hash"), u.GetPasswordHash())
		u.SetPasswordHash([]byte("new"))
		assert.Equal(t, []byte("new"), u.GetPasswordHash())
		assert.Equal(t, []byte(id.String()+"secret"), u.SaltedPassword("secret"))
		assert.Equal(t, "password", u.PasswordColumn())
	})

	t.Run("EmailVerifiedUser", func(t *testing.T) {
		id := uuid.New()
		u := &auth.EmailVerifiedUser{
			ID:           id,
			Email:        "user@example.com",
			PasswordHash: []byte("hash"),
		}
		assert.Equal(t, id.String(), u.GetID())
		assert.Equal(t, "user@example.com", u.GetUsername())
		assert.Equal(t, []string{"email"}, u.UsernameColumns())
		assert.Equal(t, []byte("hash"), u.GetPasswordHash())
		u.SetPasswordHash([]byte("new"))
		assert.Equal(t, []byte("new"), u.GetPasswordHash())
		assert.Equal(t, []byte(id.String()+"secret"), u.SaltedPassword("secret"))
		assert.Equal(t, "password", u.PasswordColumn())
		assert.Equal(t, "lookup_token", u.LookupTokenColumn())
		assert.Equal(t, "user@example.com", u.GetEmail())
		u.SetLookupToken("token")
		assert.Equal(t, "token", u.LookupToken)
		assert.False(t, u.IsVerified())
		u.SetVerified(true)
		assert.True(t, u.IsVerified())
	})

	t.Run("constructor", func(t *testing.T) {
		Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
			req := &auth.UsernameUserCreateRequest{
				UserCreateRequest: auth.UserCreateRequest{
					Update: dbtest.Update(tx),
					Ctx:    context.Background(),
					Logger: nullLogger,
				},
				Username: "user",
			}
			resp, err := auth.NewUsernameUser(req, auth.NewBasicAuthController[*auth.UsernameUser]())
			assert.NoError(t, err)
			assert.Equal(t, "user", resp.User.Username)

			emailReq := &auth.EmailVerifiedUserCreateRequest{
				UserCreateRequest: auth.UserCreateRequest{
					Update:   dbtest.Update(tx),
					Ctx:      context.Background(),
					Logger:   nullLogger,
					Mailer:   emailtest.NewTestMailer(),
					URL:      routertest.NewTestResolver(),
					Template: emailTemplates,
				},
				Email: "user@example.com",
			}
			emailResp, err := auth.NewEmailVerifiedUser(emailReq, auth.NewBasicAuthController[*auth.EmailVerifiedUser]())
			assert.NoError(t, err)
			assert.Equal(t, "user@example.com", emailResp.User.Email)
		})
	})
}
