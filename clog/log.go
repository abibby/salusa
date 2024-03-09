package clog

import (
	"context"
	"log/slog"
	"os"

	"github.com/abibby/salusa/di"
)

type key uint8

const (
	withKey key = iota
)

func Register(dp *di.DependencyProvider) {
	di.Register(dp, func(ctx context.Context, tag string) (*slog.Logger, error) {
		logger := slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
			}),
		)

		with := ctx.Value(withKey)
		if with != nil {
			logger = logger.With(with.([]any)...)
		}

		return logger, nil
	})
}

func With(ctx context.Context, attrs ...slog.Attr) context.Context {
	with := get(ctx)

	for _, attr := range attrs {
		with = append(with, attr)
	}
	return context.WithValue(ctx, withKey, with)
}

func get(ctx context.Context) []any {
	iWith := ctx.Value(withKey)
	if iWith == nil {
		return []any{}
	}
	with, ok := iWith.([]any)
	if !ok {
		return []any{}
	}
	return with
}

// func Init(ctx context.Context, logger *slog.Logger) context.Context {
// 	return context.WithValue(ctx, loggerKey, logger)
// }
// func Update(ctx context.Context, cb func(l *slog.Logger) *slog.Logger) context.Context {
// 	return context.WithValue(ctx, loggerKey, cb(Use(ctx)))
// }

// func Use(ctx context.Context) *slog.Logger {
// 	logger, ok := ctx.Value(loggerKey).(*slog.Logger)
// 	if !ok {
// 		logger = slog.Default()
// 	}
// 	return logger
// }
