package appkernel

import (
	"github.com/abibby/salusa/event/cron"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/salusadi"
	"github.com/abibby/salusa/static/template/app"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/jobs"
	"github.com/abibby/salusa/static/template/app/providers"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/database"
	"github.com/abibby/salusa/static/template/routes"
)

func Init() {
	app.Kernel = kernel.New(
		kernel.Config(config.Kernel),
		kernel.Bootstrap(
			config.Load,
			database.Init,
		),
		kernel.Providers(
			salusadi.Register,
			providers.Register,
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
}
