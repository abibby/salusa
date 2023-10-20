package jobs

import (
	"context"

	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/static/template/app/events"
)

func LogJob(ctx context.Context, e *events.LogEvent) error {
	clog.Use(ctx).Info(e.Message)
	return nil
}
