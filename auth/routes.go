package auth

import (
	"bytes"
	"context"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"reflect"
	"time"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/databasedi"
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
	Update databasedi.Update  `inject:""`
	Ctx    context.Context    `inject:""`
	URL    router.URLResolver `inject:""`
	Logger *slog.Logger       `inject:""`
}
type UserCreateResponse[T User] struct {
	User T `json:"user"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`

	Ctx  context.Context `inject:""`
	Read databasedi.Read `inject:""`
}
type LoginResponse struct {
	AccessToken  string `json:"token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh"`
	ExpiresIn    int    `json:"expires_in"`
}

type VerifyEmailRequest struct {
	Token  string             `query:"token"`
	Update databasedi.Update  `inject:""`
	Ctx    context.Context    `inject:""`
	URL    router.URLResolver `inject:""`
}

type ResetPasswordRequest struct {
	Token    string             `json:"token"`
	Password string             `json:"password"`
	Update   databasedi.Update  `inject:""`
	Ctx      context.Context    `inject:""`
	URL      router.URLResolver `inject:""`
}
type ResetPasswordResponse[T User] struct {
	User T `json:"user"`
}

type ChangePasswordRequest[T User] struct {
	OldPassword string            `json:"old_password"`
	NewPassword string            `json:"new_password"`
	User        T                 `inject:""`
	Update      databasedi.Update `inject:""`
	Ctx         context.Context   `inject:""`
}
type ChangePasswordResponse[T User] struct {
	User T `json:"user"`
}
type RefreshRequest[T User] struct {
	RefreshToken string          `json:"refresh"`
	User         T               `inject:""`
	Read         databasedi.Read `inject:""`
	Ctx          context.Context `inject:""`
}

type AuthRoutes[T User] struct {
	UserCreate     *request.RequestHandler[UserCreateRequest, *UserCreateResponse[T]]
	Login          *request.RequestHandler[LoginRequest, *LoginResponse]
	VerifyEmail    *request.RequestHandler[VerifyEmailRequest, *request.Response]
	ResetPassword  *request.RequestHandler[ResetPasswordRequest, *ResetPasswordResponse[T]]
	ChangePassword *request.RequestHandler[ChangePasswordRequest[T], *ChangePasswordResponse[T]]
	Refresh        *request.RequestHandler[RefreshRequest[T], *LoginResponse]
}

func Routes[T User](newUser func() T) *AuthRoutes[T] {

	Login := request.Handler(func(r *LoginRequest) (*LoginResponse, error) {
		u := newUser()
		var err error
		r.Read(func(tx *sqlx.Tx) error {
			u, err = builder.From[T]().
				WithContext(r.Ctx).
				Where(u.UsernameColumn(), "=", r.Username).
				First(tx)
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("failed to log in: %w", err)
		}
		if reflect.ValueOf(u).IsNil() {
			return nil, request.NewHTTPError(ErrInvalidUserPass, http.StatusUnauthorized)
		}

		if v, ok := cast[EmailVerified](u); ok {
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
			WithClaim("type", "access"),
		)
		if err != nil {
			return nil, fmt.Errorf("could not generate token: %w", err)
		}
		refresh, err := GenerateToken(
			WithSubject(u.GetID()),
			WithLifetime(time.Hour*24*30),
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

	VerifyEmail := request.Handler(func(r *VerifyEmailRequest) (*request.Response, error) {
		zeroValidated, ok := cast[EmailVerified](newUser())
		if !ok {
			return nil, request.NewHTTPError(ErrNonEmailVerifiedUser, http.StatusUnauthorized)
		}

		err := r.Update(func(tx *sqlx.Tx) error {
			u, err := builder.From[T]().
				WithContext(r.Ctx).
				Where(zeroValidated.LookupTokenColumn(), "=", r.Token).
				First(tx)
			if err != nil {
				return fmt.Errorf("failed to verify: %w", err)
			}
			if reflect.ValueOf(u).IsNil() {
				return request.NewHTTPError(ErrTokenNotFound, http.StatusUnauthorized)
			}

			v := mustCast[EmailVerified](u)
			if v.IsVerified() {
				return request.NewHTTPError(fmt.Errorf("already verified"), http.StatusUnauthorized)
			}

			v.SetVerified(true)
			v.SetLookupToken("")

			return model.SaveContext(r.Ctx, tx, u)
		})
		if err != nil {
			return nil, err
		}

		return request.Redirect(r.URL.ResolveHandler(Login)), nil
	})

	UserCreate := request.Handler(func(r *UserCreateRequest) (*UserCreateResponse[T], error) {
		u := newUser()
		u.SetUsername(r.Username)

		err := updatePassword(u, r.Password)
		if err != nil {
			return nil, err
		}

		if v, ok := cast[EmailVerified](u); ok {
			err = sendEmails(v, r.Mailer, r.URL, VerifyEmail)
			if err != nil {
				return nil, fmt.Errorf("could not send emails: %w", err)
			}
			r.Logger.Info("email verification sent", "email", v.GetEmail())
		}

		err = r.Update(func(tx *sqlx.Tx) error {
			return model.SaveContext(r.Ctx, tx, u)
		})
		if err != nil {
			return nil, err
		}

		return &UserCreateResponse[T]{
			User: u,
		}, nil
	})

	ResetPassword := request.Handler(func(r *ResetPasswordRequest) (*ResetPasswordResponse[T], error) {
		u := newUser()
		var err error
		zeroValidated, ok := cast[EmailVerified](u)
		if !ok {
			return nil, request.NewHTTPError(ErrNonEmailVerifiedUser, http.StatusUnauthorized)
		}
		r.Update(func(tx *sqlx.Tx) error {
			u, err = builder.From[T]().
				WithContext(r.Ctx).
				Where(zeroValidated.LookupTokenColumn(), "=", r.Token).
				First(tx)
			if err != nil {
				return fmt.Errorf("failed to verify: %w", err)
			}
			if reflect.ValueOf(u).IsNil() {
				return request.NewHTTPError(ErrTokenNotFound, http.StatusUnauthorized)
			}

			v := mustCast[EmailVerified](u)
			if !v.IsVerified() {
				return request.NewHTTPError(fmt.Errorf("user not verified verified"), http.StatusUnauthorized)
			}

			v.SetLookupToken("")

			err = updatePassword(u, r.Password)
			if err != nil {
				return err
			}

			return model.SaveContext(r.Ctx, tx, u)

		})
		return &ResetPasswordResponse[T]{
			User: u,
		}, nil
	})

	ChangePassword := request.Handler(func(r *ChangePasswordRequest[T]) (*ChangePasswordResponse[T], error) {
		err := bcrypt.CompareHashAndPassword(r.User.GetPasswordHash(), r.User.SaltedPassword(r.OldPassword))
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, request.NewHTTPError(ErrInvalidUserPass, http.StatusUnauthorized)
		} else if err != nil {
			return nil, fmt.Errorf("could not check password hash: %w", err)
		}

		err = updatePassword(r.User, r.NewPassword)
		if err != nil {
			return nil, err
		}

		err = r.Update(func(tx *sqlx.Tx) error {
			return model.SaveContext(r.Ctx, tx, r.User)
		})
		if err != nil {
			return nil, err
		}

		return &ChangePasswordResponse[T]{
			User: r.User,
		}, nil
	})

	Refresh := request.Handler(func(r *RefreshRequest[T]) (*LoginResponse, error) {
		claims, err := Parse(r.RefreshToken)
		if err != nil {
			return nil, err
		}

		if claims.Type != TypeRefresh {
			return nil, request.NewHTTPError(fmt.Errorf("invalid token"), http.StatusUnauthorized)
		}

		var u T
		err = r.Read(func(tx *sqlx.Tx) error {
			u, err = builder.From[T]().
				WithContext(r.Ctx).
				Find(tx, claims.Subject)
			return err
		})
		if err != nil {
			return nil, request.NewHTTPError(fmt.Errorf("failed to verify: %w", err), http.StatusUnauthorized)
		}

		expires := time.Hour

		access, err := GenerateToken(
			WithSubject(u.GetID()),
			WithLifetime(expires),
			WithIssuedAtTime(time.Now()),
			WithClaim("type", "access"),
		)
		if err != nil {
			return nil, request.NewHTTPError(fmt.Errorf("could not generate token: %w", err), http.StatusUnauthorized)
		}

		return &LoginResponse{
			AccessToken:  access,
			TokenType:    "Bearer",
			RefreshToken: r.RefreshToken,
			ExpiresIn:    int(expires.Seconds()),
		}, nil
	})

	return &AuthRoutes[T]{
		UserCreate:     UserCreate,
		Login:          Login,
		VerifyEmail:    VerifyEmail,
		ResetPassword:  ResetPassword,
		ChangePassword: ChangePassword,
		Refresh:        Refresh,
	}

}

//go:embed emails/*
var emails embed.FS

type verifyEmail struct {
	VerifyLink string
}

func updatePassword(u User, password string) error {
	hash, err := bcrypt.GenerateFromPassword(u.SaltedPassword(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.SetPasswordHash(hash)

	return nil
}

func sendEmails(v EmailVerified, mailer email.Mailer, r router.URLResolver, verifyHandler http.Handler) error {
	token := uuid.New().String()

	v.SetLookupToken(token)

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
		From:     "salusa@example.com",
		To:       []string{v.GetEmail()},
		Subject:  "Verify your account",
		HTMLBody: b.String(),
	})
	if err != nil {
		return fmt.Errorf("error sending mail: %w", err)
	}
	return nil
}

func cast[T any](v any) (T, bool) {
	t, ok := v.(T)
	return t, ok
}
func mustCast[T any](v any) T {
	return v.(T)
}

func RegisterRoutes[T User](r *router.Router, newUser func() T) {
	authRoutes := Routes(newUser)

	r.Post("/login", authRoutes.Login)
	r.Post("/login/refresh", authRoutes.Refresh)
	r.Post("/user/password/reset", authRoutes.ResetPassword)
	r.Post("/user/password/change", authRoutes.ChangePassword)
	r.Post("/user", authRoutes.UserCreate)
	r.Get("/user/verify", authRoutes.VerifyEmail)
}
