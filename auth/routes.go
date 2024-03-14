package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/abibby/salusa/request"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Ctx context.Context
	Tx  *sqlx.Tx
}
type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	Refresh   string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

var (
	ErrInvalidUserPass = errors.New("invalid username or password")
)

var Login = request.Handler(func(r *LoginRequest) (*LoginResponse, error) {
	u, err := UserQuery(r.Ctx).
		Select("id", "username", "password").
		Where("username", "=", r.Username).
		First(r.Tx)
	if err != nil {
		return nil, fmt.Errorf("failed to log in: %w", err)
	}
	if u == nil {
		return nil, request.NewHTTPError(ErrInvalidUserPass, http.StatusUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword(u.PasswordHash, u.SaltedPassword(r.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return nil, request.NewHTTPError(ErrInvalidUserPass, http.StatusUnauthorized)
	} else if err != nil {
		return nil, fmt.Errorf("could not check password hash: %w", err)
	}

	expires := time.Hour

	token, err := GenerateToken(
		WithSubject(u.ID),
		WithLifetime(expires),
		WithIssuedAtTime(time.Now()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}
	refresh, err := GenerateToken(
		WithSubject(u.ID),
		WithLifetime(time.Hour*24*30),
		WithIssuedAtTime(time.Now()),
	)
	if err != nil {
		return nil, fmt.Errorf("could not generate refresh: %w", err)
	}

	return &LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		Refresh:   refresh,
		ExpiresIn: int(expires.Seconds()),
	}, nil
})
