package auth

import (
	"bytes"
	"context"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"time"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/request"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserPass = errors.New("invalid username or password")
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Mailer email.Mailer    `inject:""`
	Tx     *sqlx.Tx        `inject:""`
	Ctx    context.Context `inject:""`
}
type UserCreateResponse[T User] struct {
	User T `json:"user"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Ctx context.Context `inject:""`
	Tx  *sqlx.Tx        `inject:""`
}
type LoginResponse struct {
	AccessToken  string `json:"token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh"`
	ExpiresIn    int    `json:"expires_in"`
}

type AuthRoutes[T User] struct {
	UserCreate *request.RequestHandler[UserCreateRequest, *UserCreateResponse[T]]
	Login      *request.RequestHandler[LoginRequest, *LoginResponse]
}

func Routes[T User](newFrom func(base *BaseUser) T) *AuthRoutes[T] {
	return &AuthRoutes[T]{
		UserCreate: request.Handler(func(r *UserCreateRequest) (*UserCreateResponse[T], error) {
			base := &BaseUser{
				ID:       uuid.New(),
				Username: r.Username,
			}

			hash, err := bcrypt.GenerateFromPassword(base.SaltedPassword(r.Password), bcrypt.MinCost)
			if err != nil {
				return nil, err
			}
			base.SetPasswordHash(hash)

			u := newFrom(base)

			var anyU any = u
			if v, ok := anyU.(EmailVerified); ok {
				err = sendEmails(v, r.Mailer)
				if err != nil {
					return nil, fmt.Errorf("could not send emails: %w", err)
				}
			}

			err = model.SaveContext(r.Ctx, r.Tx, u)
			if err != nil {
				return nil, err
			}

			return &UserCreateResponse[T]{
				User: u,
			}, nil
		}),

		Login: request.Handler(func(r *LoginRequest) (*LoginResponse, error) {
			u, err := builder.From[T]().
				WithContext(r.Ctx).
				Select("id", "username", "password").
				Where("username", "=", r.Username).
				First(r.Tx)
			if err != nil {
				return nil, fmt.Errorf("failed to log in: %w", err)
			}
			if reflect.ValueOf(u).IsNil() {
				return nil, request.NewHTTPError(ErrInvalidUserPass, http.StatusUnauthorized)
			}

			err = bcrypt.CompareHashAndPassword(u.GetPasswordHash(), u.SaltedPassword(r.Password))
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return nil, request.NewHTTPError(ErrInvalidUserPass, http.StatusUnauthorized)
			} else if err != nil {
				return nil, fmt.Errorf("could not check password hash: %w", err)
			}

			expires := time.Hour

			access, err := GenerateToken(
				WithSubject(u.GetID()),
				WithLifetime(expires),
				WithIssuedAtTime(time.Now()),
				WithClaim("type", "access"),
			)
			if err != nil {
				return nil, fmt.Errorf("could not generate token: %w", err)
			}
			refresh, err := GenerateToken(
				WithSubject(u.GetID()),
				WithLifetime(time.Hour*24*30),
				WithIssuedAtTime(time.Now()),
				WithClaim("type", "refresh"),
			)
			if err != nil {
				return nil, fmt.Errorf("could not generate refresh: %w", err)
			}

			return &LoginResponse{
				AccessToken:  access,
				TokenType:    "Bearer",
				RefreshToken: refresh,
				ExpiresIn:    int(expires.Seconds()),
			}, nil
		}),
	}

}

//go:embed emails/*
var emails embed.FS

type verifyEmail struct {
	// Email       string
	VerifyLink string
}

func sendEmails(v EmailVerified, mailer email.Mailer) error {
	token := uuid.New().String()

	v.SetValidationToken(token)

	t, err := template.ParseFS(emails, "emails/*")
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, "verify_email", &verifyEmail{
		VerifyLink: token,
	})
	if err != nil {
		return err
	}

	err = mailer.Mail(&email.Message{
		To:   []string{v.GetEmail()},
		Body: b.Bytes(),
	})
	if err != nil {
		return fmt.Errorf("error sending mail: %w", err)
	}
	return nil
}
