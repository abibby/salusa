package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/static/template/app"
	"github.com/abibby/salusa/static/template/app/appkernel"
)

func main() {
	ctx := context.Background()
	appkernel.Init()
	err := app.Kernel.Bootstrap(ctx)
	if err != nil {
		logger, err := di.Resolve[*slog.Logger](ctx, app.Kernel.DependencyProvider())
		if err != nil {
			panic(err)
		}

		logger.Error("error bootstrapping", "error", err)
		os.Exit(1)
	}

	err = app.Kernel.Run(ctx)
	if err != nil {
		logger, err := di.Resolve[*slog.Logger](ctx, app.Kernel.DependencyProvider())
		if err != nil {
			panic(err)
		}

		logger.Error("error running", "error", err)
		os.Exit(1)
	}
}
