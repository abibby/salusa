package main

import (
	"context"
	"os"

	"github.com/abibby/salusa/static/template/app"
	"github.com/abibby/salusa/static/template/app/appkernel"
)

func main() {
	ctx := context.Background()
	logger := app.Kernel.Logger(ctx)
	appkernel.Init()
	err := app.Kernel.Bootstrap(ctx)
	if err != nil {
		logger.Error("error bootstrapping", "error", err)
		os.Exit(1)
	}

	err = app.Kernel.Run(ctx)
	if err != nil {
		logger.Error("error running", "error", err)
		os.Exit(1)
	}
}
