package app

import (
	"github.com/abibby/salusa/event/cron"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/jobs"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/database"
	"github.com/abibby/salusa/static/template/routes"
)

var Kernel = kernel.NewDefaultKernel(
	kernel.Config(config.Kernel),
	kernel.Bootstrap(
		config.Load,
		database.Init,
	),
	kernel.Services(
		cron.Service().
			Schedule("* * * * *", &events.LogEvent{Message: "cron event"}),
	),
	kernel.Listeners(
		kernel.NewListener(jobs.LogJob),
	),
	kernel.InitRoutes(routes.InitRoutes),
)
