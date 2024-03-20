package routes

import (
	"fmt"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/fileserver"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/app/handlers"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/abibby/salusa/static/template/resources"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())

	auth.RegisterRoutes(r, func() *models.User {
		return &models.User{
			UsernameUser: *auth.NewBaseUser(),
		}
	})

	r.Get("/user", handlers.UserList)
	r.Get("/user/{id}", handlers.UserGet)

	r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", nil))

	fmt.Print(r.Routes())
}
