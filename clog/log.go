package clog

import (
	"context"
	"log/slog"
	"os"

	"github.com/abibby/salusa/di"
	"github.com/lmittmann/tint"
)

type key uint8

const (
	withKey key = iota
)

func Register(h slog.Handler) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		fi, err := os.Stdout.Stat()
		isTTY := err == nil && (fi.Mode()&os.ModeCharDevice) != 0

		if h == nil {
			h = tint.NewHandler(os.Stderr, &tint.Options{
				NoColor: !isTTY,
				ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
					err, ok := attr.Value.Any().(error)
					if !ok {
						return attr
					}
					errAttr := tint.Err(err)
					errAttr.Key = attr.Key
					return errAttr
				},
			})
		}
		logger := slog.New(h)

		slog.SetDefault(logger)

		di.Register(ctx, func(ctx context.Context, tag string) (*slog.Logger, error) {
			with := ctx.Value(withKey)
			if with != nil {
				return logger.With(with.([]any)...), nil
			}
			return logger, nil
		})
		return nil
	}
}

func Use(ctx context.Context) *slog.Logger {
	logger, err := di.Resolve[*slog.Logger](ctx)
	if err != nil {
		return slog.Default()
	}
	return logger
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
