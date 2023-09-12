package routes

import (
	"github.com/abibby/salusa/fileserver"
	"github.com/abibby/salusa/router"
	"github.com/abibby/salusa/static/template/app/controllers"
	"github.com/abibby/salusa/static/template/resources"
)

func InitRoutes(r *router.Router) {
	r.Get("/add", controllers.Add)

	r.Handle("/", fileserver.WithFallback(resources.Content, "dist", "index.html", nil))
}
