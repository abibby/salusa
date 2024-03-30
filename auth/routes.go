package auth

import (
	"context"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/email"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/view"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserPass      = errors.New("invalid username or password")
	ErrTokenNotFound        = errors.New("token not found")
	ErrNonEmailVerifiedUser = errors.New("non email verified user")
)

type AuthRoutes[T User] struct {
	UserCreate     *request.RequestHandler[UserCreateRequest, *UserCreateResponse[T]]
	Login          *request.RequestHandler[LoginRequest, *LoginResponse]
	VerifyEmail    *request.RequestHandler[VerifyEmailRequest, http.Handler]
	ResetPassword  *request.RequestHandler[ResetPasswordRequest, *ResetPasswordResponse[T]]
	ForgotPassword *request.RequestHandler[ForgotPasswordRequest, *ForgotPasswordResponse]
	ChangePassword *request.RequestHandler[ChangePasswordRequest[T], *ChangePasswordResponse[T]]
	Refresh        *request.RequestHandler[RefreshRequest[T], *LoginResponse]
}

type RouteOptions[T User, R any] struct {
	NewUser           func(request R) T
	ResetPasswordName string
}

//go:embed emails/*
var emails embed.FS
var defaultViewTemplate = view.NewViewTemplate(emails, "**/*.html")

func Routes[T User, R any](newUser func(request R) T, resetPasswordName string) *AuthRoutes[T] {
	options := &RouteOptions[T, R]{
		NewUser:           newUser,
		ResetPasswordName: resetPasswordName,
	}
	return &AuthRoutes[T]{
		UserCreate:     options.userCreate(),
		Login:          options.login(),
		VerifyEmail:    options.verifyEmail(),
		ResetPassword:  options.resetPassword(),
		ChangePassword: options.changePassword(),
		Refresh:        options.refresh(),
		ForgotPassword: options.forgotPassword(),
	}

}

func RegisterRoutes[T User, R any](r *router.Router, newUser func(request R) T, resetPasswordName string) {
	authRoutes := Routes(newUser, resetPasswordName)

	r.Post("/login", authRoutes.Login).Name("auth.login")
	r.Post("/user/password/reset", authRoutes.ResetPassword).Name("auth.password.reset")
	r.Post("/user/password/forgot", authRoutes.ForgotPassword).Name("auth.password.forgot")
	r.Post("/user", authRoutes.UserCreate).Name("auth.user.create")
	r.Get("/user/verify", authRoutes.VerifyEmail).Name("auth.email.verify")
	r.Post("/login/refresh", authRoutes.Refresh).Name("auth.refresh")

	r.Group("", func(r *router.Router) {
		r.Use(AttachUser())
		r.Use(LoggedIn())

		r.Post("/user/password/change", authRoutes.ChangePassword).Name("auth.password.change")
	})
}

type UserCreateRequest struct {
	Password string `json:"password"`

	Mailer   email.Mailer       `inject:""`
	Update   databasedi.Update  `inject:""`
	Ctx      context.Context    `inject:""`
	Logger   *slog.Logger       `inject:""`
	Request  *http.Request      `inject:""`
	URL      router.URLResolver `inject:""`
	Template *view.ViewTemplate `inject:",optional"`
}
type UserCreateResponse[T User] struct {
	User T `json:"user"`
}

func (o *RouteOptions[T, R]) userCreate() *request.RequestHandler[UserCreateRequest, *UserCreateResponse[T]] {
	return request.Handler(func(r *UserCreateRequest) (*UserCreateResponse[T], error) {
		req, err := helpers.NewOf[R]()
		if err != nil {
			return nil, err
		}
		err = request.Run(r.Request, req)
		if err != nil {
			return nil, err
		}
		u := o.NewUser(req)

		for _, col := range u.UsernameColumns() {
			v, err := helpers.RGetValue(reflect.ValueOf(u), col)
			if err != nil {
				return nil, err
			}
			if v.Kind() == reflect.String {
				v.Set(reflect.ValueOf(strings.ToLower(v.String())))
			}
		}

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
				go func() {
					err = o.sendEmails(&sendEmailOptions{
						URL:          r.URL,
						ViewTemplate: r.Template,
						User:         v,
						Mailer:       r.Mailer,
						TemplateName: "verify_email.html",
						Subject:      "Verify your email",
					})
					if err != nil {
						r.Logger.Info("could not send verification email", "email", v.GetEmail(), "error", err)
						return
					}
					r.Logger.Info("email verification sent", "email", v.GetEmail())
				}()
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

func (o *RouteOptions[T, R]) login() *request.RequestHandler[LoginRequest, *LoginResponse] {
	return request.Handler(func(r *LoginRequest) (*LoginResponse, error) {
		u, err := helpers.NewOf[T]()
		if err != nil {
			return nil, err
		}
		userColumns := u.UsernameColumns()
		if len(userColumns) == 0 {
			panic("need columns")
		}
		err = r.Read(func(tx *sqlx.Tx) error {
			q := builder.From[T]().WithContext(r.Ctx)
			for _, column := range userColumns {
				q = q.OrWhere(column, "=", strings.ToLower(r.Username))
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

type VerifyEmailRequest struct {
	Token  string             `query:"token"`
	Update databasedi.Update  `inject:""`
	Ctx    context.Context    `inject:""`
	URL    router.URLResolver `inject:""`
}

func (o *RouteOptions[T, R]) verifyEmail() *request.RequestHandler[VerifyEmailRequest, http.Handler] {
	return request.Handler(func(r *VerifyEmailRequest) (http.Handler, error) {
		v, err := helpers.NewOf[T]()
		if err != nil {
			return nil, err
		}
		emptyValidated, ok := cast[EmailVerified](v)
		if !ok {
			return nil, request.NewHTTPError(ErrNonEmailVerifiedUser, http.StatusUnauthorized)
		}

		err = r.Update(func(tx *sqlx.Tx) error {
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

type ResetPasswordRequest struct {
	Token    string             `json:"token" validate:"required|min:1"`
	Password string             `json:"password" validate:"required"`
	Update   databasedi.Update  `inject:""`
	Ctx      context.Context    `inject:""`
	URL      router.URLResolver `inject:""`
}
type ResetPasswordResponse[T User] struct {
	User T `json:"user"`
}

func (o *RouteOptions[T, R]) resetPassword() *request.RequestHandler[ResetPasswordRequest, *ResetPasswordResponse[T]] {
	return request.Handler(func(r *ResetPasswordRequest) (*ResetPasswordResponse[T], error) {
		u, err := helpers.NewOf[T]()
		if err != nil {
			return nil, err
		}
		zeroValidated, ok := cast[EmailVerified](u)
		if !ok {
			return nil, ErrNonEmailVerifiedUser
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

type ForgotPasswordRequest struct {
	Email    string             `json:"email" validate:"required|email"`
	Update   databasedi.Update  `inject:""`
	Ctx      context.Context    `inject:""`
	Mailer   email.Mailer       `inject:""`
	Logger   *slog.Logger       `inject:""`
	URL      router.URLResolver `inject:""`
	Template *view.ViewTemplate `inject:",optional"`
}
type ForgotPasswordResponse struct {
}

func (o *RouteOptions[T, R]) forgotPassword() *request.RequestHandler[ForgotPasswordRequest, *ForgotPasswordResponse] {
	return request.Handler(func(r *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
		u, err := helpers.NewOf[T]()
		if err != nil {
			return nil, err
		}
		_, ok := cast[EmailVerified](u)
		if !ok {
			panic("not email verified")
		}
		userColumns := u.UsernameColumns()
		if len(userColumns) == 0 {
			panic("need columns")
		}
		err = r.Update(func(tx *sqlx.Tx) error {
			q := builder.From[T]().WithContext(r.Ctx)
			for _, column := range userColumns {
				q = q.OrWhere(column, "=", strings.ToLower(r.Email))
			}

			u, err := q.First(tx)
			if err != nil {
				return err
			}
			if reflect.ValueOf(u).IsZero() {
				r.Logger.Info("password reset attempt for unused email", slog.String("email", r.Email))
				return nil
			}

			go func() {
				err = o.sendEmails(&sendEmailOptions{
					URL:          r.URL,
					ViewTemplate: r.Template,
					User:         mustCast[EmailVerified](u),
					Mailer:       r.Mailer,
					TemplateName: "reset_password.html",
					Subject:      "Password reset",
				})
				if err != nil {
					r.Logger.Info("password failed to send", slog.String("email", r.Email), slog.Any("error", err))
					return
				}
				r.Logger.Info("password reset sent", slog.String("email", r.Email))
			}()
			return model.SaveContext(r.Ctx, tx, u)
		})
		if err != nil {
			return nil, err
		}
		return &ForgotPasswordResponse{}, err
	})
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

func (o *RouteOptions[T, R]) changePassword() *request.RequestHandler[ChangePasswordRequest[T], *ChangePasswordResponse[T]] {
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

type RefreshRequest[T User] struct {
	RefreshToken string          `json:"refresh"`
	User         T               `inject:""`
	Read         databasedi.Read `inject:""`
	Ctx          context.Context `inject:""`
}

func (o *RouteOptions[T, R]) refresh() *request.RequestHandler[RefreshRequest[T], *LoginResponse] {
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

func updatePassword(u User, password string) error {
	hash, err := bcrypt.GenerateFromPassword(u.SaltedPassword(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.SetPasswordHash(hash)

	return nil
}

type sendEmailOptions struct {
	URL          router.URLResolver
	ViewTemplate *view.ViewTemplate
	User         EmailVerified
	Mailer       email.Mailer
	TemplateName string
	Subject      string
}

func (o *RouteOptions[T, R]) sendEmails(opt *sendEmailOptions) error {
	token := uuid.New().String()

	opt.User.SetLookupToken(token)

	if opt.ViewTemplate == nil {
		opt.ViewTemplate = defaultViewTemplate
	}

	viewResponse, err := view.View(opt.TemplateName, map[string]any{
		"ResetPasswordName": o.ResetPasswordName,
		"Token":             token,
	}).Run(&view.ViewRequest{
		URL:      opt.URL,
		Template: opt.ViewTemplate,
	})
	if err != nil {
		return err
	}
	b, err := viewResponse.Bytes()
	if err != nil {
		return err
	}
	err = opt.Mailer.Mail(&email.Message{
		From:     "salusa@example.com",
		To:       []string{opt.User.GetEmail()},
		Subject:  opt.Subject,
		HTMLBody: string(b),
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
