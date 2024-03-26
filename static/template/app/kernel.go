package app

import (
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/event/cron"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/salusadi"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/jobs"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/abibby/salusa/static/template/app/providers"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/database"
	"github.com/abibby/salusa/static/template/resources"
	"github.com/abibby/salusa/static/template/routes"
	"github.com/abibby/salusa/view"
)

var Kernel = kernel.New[*config.Config](
	kernel.Config(config.Load),
	kernel.Bootstrap(
		salusadi.Register[*models.User],
		event.RegisterChannelQueue,
		providers.Register,
		database.Init,
		view.Register(resources.Content, "**/*.html"),
	),
	kernel.Services(
		cron.Service().
			Schedule("* * * * *", &events.LogEvent{Message: "cron event"}),
		event.Service(
			event.NewListener[*jobs.LogJob](),
		),
	),
	kernel.InitRoutes(routes.InitRoutes),
)
