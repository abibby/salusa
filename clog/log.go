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

func Register(h slog.Handler) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		di.Register(ctx, func(ctx context.Context, tag string) (*slog.Logger, error) {
			if h == nil {
				h = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
					AddSource: true,
				})
			}
			logger := slog.New(h)

			with := ctx.Value(withKey)
			if with != nil {
				logger = logger.With(with.([]any)...)
			}

			return logger, nil
		})
		return nil
	}
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
