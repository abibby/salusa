package app

import (
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/database"
	"github.com/abibby/salusa/static/template/routes"
)

var Kernel = kernel.New(
	kernel.Bootstrap(
		config.Load,
		database.Init,
	),
	kernel.InitRoutes(routes.InitRoutes),
	kernel.Middleware(
		request.HandleErrors(),
		request.WithDB(database.DB),
	),
)
