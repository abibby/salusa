package routes

import (
	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/fileserver"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/abibby/salusa/static/template/resources"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	// r.Use(databasedi.Middleware())

	authRoutes := auth.Routes(func() *models.User {
		return &models.User{
			UsernameUser: *auth.NewBaseUser(),
		}
	})

	r.Get("login", authRoutes.Login)
	r.Post("user", authRoutes.UserCreate)

	// r.Get("/user", handlers.UserList)
	// r.Get("/user/{id}", handlers.UserGet)
	// r.Post("/user", handlers.UserCreate)

	r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", nil))
}
