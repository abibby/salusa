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

func register() (*di.DependencyProvider, *bytes.Buffer) {
	dp := di.NewDependencyProvider()
	b := bytes.NewBuffer([]byte{})

	clog.Register(dp, slog.NewTextHandler(b, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Time(slog.TimeKey, time.Time{})
			}
			return a
		},
	}))

	return dp, b
}

func TestWith(t *testing.T) {
	t.Run("with string", func(t *testing.T) {
		dp, b := register()
		ctx := context.Background()

		ctx = clog.With(ctx, slog.String("foo", "bar"))

		l, err := di.Resolve[*slog.Logger](ctx, dp)
		assert.NoError(t, err)

		l.Warn("test")
		assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=WARN msg=test foo=bar\n", b.String())
	})

	t.Run("with multiple", func(t *testing.T) {
		dp, b := register()
		ctx := context.Background()

		ctx = clog.With(ctx, slog.String("a", "1"))
		ctx = clog.With(ctx, slog.String("b", "2"))

		l, err := di.Resolve[*slog.Logger](ctx, dp)
		assert.NoError(t, err)

		l.Warn("test")
		assert.Equal(t, "time=0001-01-01T00:00:00.000Z level=WARN msg=test a=1 b=2\n", b.String())
	})
}

func TestResolve(t *testing.T) {
	t.Run("no handler", func(t *testing.T) {
		dp := di.NewDependencyProvider()

		clog.Register(dp, nil)

		ctx := context.Background()

		l, err := di.Resolve[*slog.Logger](ctx, dp)
		assert.NoError(t, err)
		assert.NotNil(t, l)
	})
}
