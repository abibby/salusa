package main

import (
	"context"
	"os"

	"github.com/abibby/salusa/static/template/app"
)

func main() {
	ctx := context.Background()

	err := app.Kernel.Bootstrap(ctx)
	if err != nil {
		app.Kernel.Logger(ctx).Error("error bootstrapping", "error", err)
		os.Exit(1)
	}

	err = app.Kernel.Run(ctx)
	if err != nil {
		app.Kernel.Logger(ctx).Error("error running", "error", err)
		os.Exit(1)
	}
}
