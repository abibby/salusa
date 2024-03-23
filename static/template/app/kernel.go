package app

import (
	"context"

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
		database.Init,
		func(ctx context.Context, k *kernel.Kernel) error {
			resources.View = view.Factory(k.DependencyProvider(), resources.Content, "**/*.html")
			return nil
		},
	),
	kernel.Providers(
		salusadi.Register[*models.User],
		event.RegisterChannelQueue,
		providers.Register,
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
