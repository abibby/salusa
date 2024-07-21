package app

import (
	"context"

	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/event/cron"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/openapidoc"
	"github.com/abibby/salusa/salusadi"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/jobs"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/abibby/salusa/static/template/app/providers"
	"github.com/abibby/salusa/static/template/config"
	"github.com/abibby/salusa/static/template/migrations"
	"github.com/abibby/salusa/static/template/resources"
	"github.com/abibby/salusa/static/template/routes"
	"github.com/abibby/salusa/view"
	"github.com/go-openapi/spec"
	"github.com/google/uuid"
)

var Kernel = kernel.New(
	kernel.Config(config.Load),
	kernel.Bootstrap(
		salusadi.Register[*models.User](migrations.Use()),
		view.Register(resources.Content, "**/*.html"),
		providers.Register,
		func(ctx context.Context) error {
			openapidoc.RegisterFormat[uuid.UUID]("uuid")
			return nil
		},
	),
	kernel.Services(
		cron.Service().
			Schedule("* * * * *", &events.LogEvent{Message: "cron event"}),
		event.Service(
			event.NewListener[*jobs.LogJob](),
		),
	),
	kernel.InitRoutes(routes.InitRoutes),
	kernel.APIDocumentationInfo(spec.InfoProps{
		Title:       "Salusa Example API",
		Description: `This is the API documentaion for the example Salusa application`,
	}),
	kernel.APIDocumentationBasePath("/api"),
)
