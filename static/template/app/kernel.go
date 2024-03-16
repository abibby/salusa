package app

import (
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/event/cron"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/salusadi"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/jobs"
	"github.com/abibby/salusa/static/template/app/providers"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/database"
	"github.com/abibby/salusa/static/template/routes"
)

var Kernel = kernel.New(
	kernel.Config(config.Kernel),
	kernel.Bootstrap(
		config.Load,
		database.Init,
	),
	kernel.Providers(
		salusadi.Register,
		providers.Register,
		event.RegisterChannelQueue,
	),
	kernel.Services(
		cron.Service().
			Schedule("* * * * *", &events.LogEvent{Message: "cron event"}),
		event.Service().
			Add(event.NewListener[*jobs.LogJob]()),
	),
	kernel.InitRoutes(routes.InitRoutes),
)
