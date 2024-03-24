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
	"github.com/abibby/salusa/internal/helpers"
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
	Password string `json:"password"`

	Mailer  email.Mailer       `inject:""`
	Update  databasedi.Update  `inject:""`
	Ctx     context.Context    `inject:""`
	URL     router.URLResolver `inject:""`
	Logger  *slog.Logger       `inject:""`
	Request *http.Request      `inject:""`
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
	VerifyEmail    *request.RequestHandler[VerifyEmailRequest, http.Handler]
	ResetPassword  *request.RequestHandler[ResetPasswordRequest, *ResetPasswordResponse[T]]
	ChangePassword *request.RequestHandler[ChangePasswordRequest[T], *ChangePasswordResponse[T]]
	Refresh        *request.RequestHandler[RefreshRequest[T], *LoginResponse]
}

func Routes[T User, R any](newUser func(request R) T) *AuthRoutes[T] {
	verify := VerifyEmail[T]()

	return &AuthRoutes[T]{
		UserCreate:     UserCreate(newUser, verify),
		Login:          Login[T](),
		VerifyEmail:    verify,
		ResetPassword:  ResetPassword[T](),
		ChangePassword: ChangePassword[T](),
		Refresh:        Refresh[T](),
	}

}

func UserCreate[T User, R any](newUser func(request R) T, verifyEmail http.Handler) *request.RequestHandler[UserCreateRequest, *UserCreateResponse[T]] {
	return request.Handler(func(r *UserCreateRequest) (*UserCreateResponse[T], error) {
		req := helpers.NewOf[R]()
		err := request.Run(r.Request, req)
		if err != nil {
			return nil, err
		}
		u := newUser(req)

		err = r.Update(func(tx *sqlx.Tx) error {
			err := model.SaveContext(r.Ctx, tx, u)
			if err != nil {
				return err
			}

			err = updatePassword(u, r.Password)
			if err != nil {
				return err
			}

			if v, ok := cast[EmailVerified](u); ok {
				r.Logger.Info("email verification sent", "email", v.GetEmail())
				err = sendEmails(v, r.Mailer, r.URL, verifyEmail)
				if err != nil {
					return fmt.Errorf("could not send emails: %w", err)
				}
				r.Logger.Info("email verification sent", "email", v.GetEmail())
			}

			return model.SaveContext(r.Ctx, tx, u)
		})
		if err != nil {
			return nil, err
		}

		return &UserCreateResponse[T]{
			User: u,
		}, nil
	})
}

func Login[T User]() *request.RequestHandler[LoginRequest, *LoginResponse] {
	return request.Handler(func(r *LoginRequest) (*LoginResponse, error) {
		u := helpers.NewOf[T]()
		userColumns := u.UsernameColumns()
		if len(userColumns) == 0 {
			panic("need columns")
		}
		err := r.Read(func(tx *sqlx.Tx) error {
			q := builder.From[T]().WithContext(r.Ctx)
			for _, column := range userColumns {
				q = q.OrWhere(column, "=", r.Username)
			}

			var err error
			u, err = q.First(tx)
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
				return nil, Err401Unauthorized
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
}

func VerifyEmail[T User]() *request.RequestHandler[VerifyEmailRequest, http.Handler] {
	return request.Handler(func(r *VerifyEmailRequest) (http.Handler, error) {
		emptyValidated, ok := cast[EmailVerified](helpers.NewOf[T]())
		if !ok {
			return nil, request.NewHTTPError(ErrNonEmailVerifiedUser, http.StatusUnauthorized)
		}

		err := r.Update(func(tx *sqlx.Tx) error {
			u, err := builder.From[T]().
				WithContext(r.Ctx).
				Where(emptyValidated.LookupTokenColumn(), "=", r.Token).
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

		return http.RedirectHandler(r.URL.Resolve("login"), http.StatusFound), nil
	})
}

func ResetPassword[T User]() *request.RequestHandler[ResetPasswordRequest, *ResetPasswordResponse[T]] {
	return request.Handler(func(r *ResetPasswordRequest) (*ResetPasswordResponse[T], error) {
		u := helpers.NewOf[T]()
		var err error
		zeroValidated, ok := cast[EmailVerified](u)
		if !ok {
			return nil, request.NewHTTPError(ErrNonEmailVerifiedUser, http.StatusUnauthorized)
		}
		err = r.Update(func(tx *sqlx.Tx) error {
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
				return Err401Unauthorized
			}

			v.SetLookupToken("")

			err = updatePassword(u, r.Password)
			if err != nil {
				return err
			}

			return model.SaveContext(r.Ctx, tx, u)

		})
		if err != nil {
			return nil, err
		}
		return &ResetPasswordResponse[T]{
			User: u,
		}, nil
	})
}

func ChangePassword[T User]() *request.RequestHandler[ChangePasswordRequest[T], *ChangePasswordResponse[T]] {
	return request.Handler(func(r *ChangePasswordRequest[T]) (*ChangePasswordResponse[T], error) {
		err := bcrypt.CompareHashAndPassword(r.User.GetPasswordHash(), r.User.SaltedPassword(r.OldPassword))
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, Err401Unauthorized
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
}

func Refresh[T User]() *request.RequestHandler[RefreshRequest[T], *LoginResponse] {
	return request.Handler(func(r *RefreshRequest[T]) (*LoginResponse, error) {
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

func RegisterRoutes[T User, R any](r *router.Router, newUser func(request R) T) {
	authRoutes := Routes(newUser)

	r.Post("/login", authRoutes.Login).Name("auth.login")
	r.Post("/login/refresh", authRoutes.Refresh).Name("auth.refresh")
	r.Post("/user/password/reset", authRoutes.ResetPassword).Name("auth.password.reset")
	r.Post("/user/password/change", authRoutes.ChangePassword).Name("auth.password.change")
	r.Post("/user", authRoutes.UserCreate).Name("auth.user.create")
	r.Get("/user/verify", authRoutes.VerifyEmail).Name("auth.email.verify")
}
