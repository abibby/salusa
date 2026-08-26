package app

import (
	"context"

	"abibby.com/salusa/auth"
	"abibby.com/salusa/clog"
	"abibby.com/salusa/database"
	"abibby.com/salusa/email"
	"abibby.com/salusa/event"
	"abibby.com/salusa/event/cron"
	"abibby.com/salusa/filesystem"
	"abibby.com/salusa/kernel"
	"abibby.com/salusa/openapidoc"
	"abibby.com/salusa/openapidoc/openapidocdi"
	"abibby.com/salusa/pubsub/channelpubsub"
	"abibby.com/salusa/request"
	"abibby.com/salusa/static/template/app/events"
	"abibby.com/salusa/static/template/app/jobs"
	"abibby.com/salusa/static/template/app/models"
	"abibby.com/salusa/static/template/app/providers"
	"abibby.com/salusa/static/template/config"
	"abibby.com/salusa/static/template/migrations"
	"abibby.com/salusa/static/template/resources"
	"abibby.com/salusa/static/template/routes"
	"abibby.com/salusa/view"
	"github.com/go-openapi/spec"
)

var Kernel = kernel.New(
	kernel.Config(config.Load),
	kernel.Bootstrap(
		view.Register(resources.Content, "**/*.html"),
		providers.Register,
		kernel.Register(func(ctx context.Context, c *config.Config) {
			database.Register(ctx, c.Database, migrations.Use())
			email.Register(ctx, c.Mail)
			channelpubsub.Register(ctx)

			clog.RegisterDefault(ctx)
			request.Register(ctx)
			auth.Register[*models.User](ctx)
			event.Register(ctx)
			filesystem.Register(ctx)
			openapidocdi.Register(ctx)
		}),
	),
	kernel.Services(
		cron.Service().
			Schedule("* * * * *", &events.LogEvent{Message: "cron event"}),
		event.Service(
			event.NewListener[*jobs.LogJob](),
		),
	),
	kernel.InitRoutes(routes.InitRoutes),
	kernel.APIDocumentation(
		openapidoc.Info(spec.InfoProps{
			Title:       "Salusa Example API",
			Description: `This is the API documentaion for the example Salusa application`,
		}),
		openapidoc.BasePath("/api"),
	),
)
