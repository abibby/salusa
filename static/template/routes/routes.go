package routes

import (
	"github.com/abibby/salusa/database/databasedi"
	"github.com/abibby/salusa/fileserver"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/app/handlers"
	"github.com/abibby/salusa/static/template/resources"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	r.Use(databasedi.Middleware())

	r.Get("/user", handlers.UserList)
	r.Get("/user/{id}", handlers.UserGet)
	r.Post("/user", handlers.UserCreate)

	r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", nil))
}
