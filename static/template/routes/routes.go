package routes

import (
	"net/http"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/openapidoc"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/app/handlers"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/abibby/salusa/view"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	r.Use(auth.AttachUser())

	auth.RegisterRoutes(r, auth.NewBasicAuthController[*models.User](
		auth.NewUser(func(r *auth.EmailVerifiedUserCreateRequest) *models.User {
			return &models.User{
				EmailVerifiedUser: *auth.NewEmailVerifiedUser(r),
			}
		}),
		auth.ResetPasswordName("reset-password"),
	))

	r.Get("/", view.View("index.html", nil)).Name("home")
	r.Get("/login", view.View("login.html", nil)).Name("login")
	r.Get("/user/create", view.View("create_user.html", nil)).Name("user.create")

	r.Get("/user", handlers.UserList)
	r.Get("/user/{id}", handlers.UserGet)

	r.Get("/swagger.json", handlers.Docs)

	r.Handle("/docs", http.StripPrefix("/docs", openapidoc.SwaggerUI("/swagger.json")))
}
