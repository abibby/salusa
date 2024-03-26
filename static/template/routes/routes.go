package routes

import (
	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/app/handlers"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/abibby/salusa/view"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	r.Use(auth.AttachUser())

	auth.RegisterRoutes(r, func(r *auth.EmailVerifiedUserCreateRequest) *models.User {
		return &models.User{
			EmailVerifiedUser: *auth.NewEmailVerifiedUser(r),
		}
	}, "/reset-password")

	r.Get("/", view.View("index.html", nil)).Name("home")
	r.Get("/login", view.View("login.html", nil)).Name("login")
	r.Get("/user/create", view.View("create_user.html", nil)).Name("user.create")

	r.Get("/user", handlers.UserList)
	r.Get("/user/{id}", handlers.UserGet)

	// r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", nil))

	r.PrintRoutes()
}
