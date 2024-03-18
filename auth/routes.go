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
	"github.com/abibby/salusa/router"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserPass      = errors.New("invalid username or password")
	ErrTokenNotFound        = errors.New("token not found")
	ErrNonEmailVerifiedUser = errors.New("non email verified user")
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Mailer email.Mailer       `inject:""`
	Tx     *sqlx.Tx           `inject:""`
	Ctx    context.Context    `inject:""`
	URL    router.URLResolver `inject:""`
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

type VerifyEmailRequest struct {
	Token string             `query:"token"`
	Tx    *sqlx.Tx           `inject:""`
	Ctx   context.Context    `inject:""`
	URL   router.URLResolver `inject:""`
}

type AuthRoutes[T User] struct {
	UserCreate  *request.RequestHandler[UserCreateRequest, *UserCreateResponse[T]]
	Login       *request.RequestHandler[LoginRequest, *LoginResponse]
	VerifyEmail *request.RequestHandler[VerifyEmailRequest, *request.Response]
}

func Routes[T User](newUser func() T) *AuthRoutes[T] {
	VerifyEmail := request.Handler(func(r *VerifyEmailRequest) (*request.Response, error) {
		var zeroUser T
		var u model.Model = zeroUser
		_, ok := u.(EmailVerified)
		if !ok {
			return nil, request.NewHTTPError(ErrNonEmailVerifiedUser, http.StatusUnauthorized)
		}
		u, err := builder.From[T]().
			WithContext(r.Ctx).
			Select("id", zeroUser.UsernameColumn(), zeroUser.PasswordColumn()).
			Where("verification_token", "=", r.Token).
			First(r.Tx)
		if err != nil {
			return nil, fmt.Errorf("failed to verify: %w", err)
		}
		if reflect.ValueOf(u).IsNil() {
			return nil, request.NewHTTPError(ErrTokenNotFound, http.StatusUnauthorized)
		}

		v := u.(EmailVerified)
		v.SetVerified(true)
		v.SetVerificationToken("")

		err = model.SaveContext(r.Ctx, r.Tx, u)
		if err != nil {
			return nil, err
		}

		return request.Redirect(r.URL.Resolve("/")), nil
	})

	UserCreate := request.Handler(func(r *UserCreateRequest) (*UserCreateResponse[T], error) {
		u := newUser()
		u.SetUsername(r.Username)

		hash, err := bcrypt.GenerateFromPassword(u.SaltedPassword(r.Password), bcrypt.MinCost)
		if err != nil {
			return nil, err
		}
		u.SetPasswordHash(hash)

		var anyUser any = u
		if v, ok := anyUser.(EmailVerified); ok {
			err = sendEmails(v, r.Mailer, r.URL, VerifyEmail)
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
	})

	Login := request.Handler(func(r *LoginRequest) (*LoginResponse, error) {
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

		var anyUser any = u
		if v, ok := anyUser.(EmailVerified); ok {
			if !v.IsVerified() {
				return nil, request.NewHTTPError(ErrInvalidUserPass, http.StatusUnauthorized)
			}
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
	})
	return &AuthRoutes[T]{
		UserCreate:  UserCreate,
		Login:       Login,
		VerifyEmail: VerifyEmail,
	}

}

//go:embed emails/*
var emails embed.FS

type verifyEmail struct {
	// Email       string
	VerifyLink string
}

func sendEmails(v EmailVerified, mailer email.Mailer, r router.URLResolver, verifyHandler http.Handler) error {
	token := uuid.New().String()

	v.SetVerificationToken(token)

	t, err := template.ParseFS(emails, "emails/*")
	if err != nil {
		return err
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, "verify_email.html", &verifyEmail{
		VerifyLink: r.ResolveHandler(verifyHandler, "token", token),
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
