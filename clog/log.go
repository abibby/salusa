package clog

import (
	"context"

	"log/slog"
)

type key string

var loggerKey = key("logger")

func Init(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
func Update(ctx context.Context, cb func(l *slog.Logger) *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, cb(Use(ctx)))
}

func Use(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
	if !ok {
		logger = slog.Default()
	}
	return logger
}
