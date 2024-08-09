package clog

import (
	"context"
	"log/slog"
	"os"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusaconfig"
	"github.com/lmittmann/tint"
)

type key uint8

const (
	withKey key = iota
)

type rootLogger slog.Logger

type LoggerConfiger interface {
	LoggerConfig() Config
}

type Config interface {
	Handler() (slog.Handler, error)
}

type DefaultConfig struct{}

var _ Config = (*DefaultConfig)(nil)

func NewDefaultConfig() *DefaultConfig {
	return &DefaultConfig{}
}

func (c *DefaultConfig) Handler() (slog.Handler, error) {
	return DefaultHandler(), nil
}

func Register(ctx context.Context) error {
	di.RegisterLazySingletonWith(ctx, func(cfg salusaconfig.Config) (*rootLogger, error) {
		var h slog.Handler
		if lc, ok := cfg.(LoggerConfiger); ok {
			var err error
			h, err = lc.LoggerConfig().Handler()
			if err != nil {
				return nil, err
			}
		}
		if h == nil {
			h = DefaultHandler()
		}
		logger := slog.New(h)

		slog.SetDefault(logger)
		return (*rootLogger)(logger), nil
	})

	registerLogger(ctx)
	return nil
}
func RegisterWith(h slog.Handler) func(ctx context.Context) error {
	return func(ctx context.Context) error {

		di.RegisterLazySingleton(ctx, func() (*rootLogger, error) {
			if h == nil {
				h = DefaultHandler()
			}
			logger := slog.New(h)

			slog.SetDefault(logger)
			return (*rootLogger)(logger), nil
		})

		registerLogger(ctx)

		return nil
	}
}
func DefaultHandler() slog.Handler {
	fi, err := os.Stdout.Stat()
	isTTY := err == nil && (fi.Mode()&os.ModeCharDevice) != 0
	return tint.NewHandler(os.Stderr, &tint.Options{
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
func registerLogger(ctx context.Context) {
	di.RegisterWith(ctx, func(ctx context.Context, tag string, logger *rootLogger) (*slog.Logger, error) {
		with := ctx.Value(withKey)
		sLogger := (*slog.Logger)(logger)
		if with != nil {
			return sLogger.With(with.([]any)...), nil
		}
		return sLogger, nil
	})
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
