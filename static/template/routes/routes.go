package routes

import (
	"github.com/abibby/salusa/fileserver"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/app/handlers"
	"github.com/abibby/salusa/static/template/database"
	"github.com/abibby/salusa/static/template/resources"
)

func InitRoutes(r *router.Router) {
	r.Use(request.HandleErrors())
	r.Use(request.WithDB(database.DB))

	r.Get("/add", handlers.Add)

	r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", nil))
}
