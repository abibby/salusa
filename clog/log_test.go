package clog_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/di"
	"github.com/stretchr/testify/assert"
)

func register() (context.Context, *bytes.Buffer) {
	ctx := di.ContextWithDependencyProvider(
		context.Background(),
		di.NewDependencyProvider(),
	)

	b := bytes.NewBuffer([]byte{})

	_ = clog.RegisterWith(slog.NewTextHandler(b, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Time(slog.TimeKey, time.Time{})
			}
			return a
		},
	}))(ctx)

	return ctx, b
}

func TestWith(t *testing.T) {
	t.Run("with string", func(t *testing.T) {
		ctx, b := register()

		ctx = clog.With(ctx, slog.String("foo", "bar"))

		l, err := di.Resolve[*slog.Logger](ctx)
		assert.NoError(t, err)

		l.Warn("test")
		assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=WARN msg=test foo=bar\n", b.String())
	})

	t.Run("with multiple", func(t *testing.T) {
		ctx, b := register()

		ctx = clog.With(ctx, slog.String("a", "1"))
		ctx = clog.With(ctx, slog.String("b", "2"))

		l, err := di.Resolve[*slog.Logger](ctx)
		assert.NoError(t, err)

		l.Warn("test")
		assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=WARN msg=test a=1 b=2\n", b.String())
	})
}

func TestResolve(t *testing.T) {
	t.Run("no handler", func(t *testing.T) {
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)

		_ = clog.RegisterWith(nil)(ctx)

		l, err := di.Resolve[*slog.Logger](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, l)
	})
}
