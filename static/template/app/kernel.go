package app

import (
	"github.com/abibby/salusa/event/cron"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/jobs"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/database"
	"github.com/abibby/salusa/static/template/routes"
)

var Kernel = kernel.NewDefaultKernel(
	kernel.Port(6900),
	kernel.Bootstrap(
		config.Load,
		database.Init,
	),
	kernel.Services(
		cron.NewService().Schedule("* * * * *", &events.LogEvent{}),
	),
	kernel.Listeners(
		kernel.NewListener(jobs.LogJob),
	),
	kernel.InitRoutes(routes.InitRoutes),
	kernel.Middleware(
		request.HandleErrors(),
		request.WithDB(database.DB),
	),
)
