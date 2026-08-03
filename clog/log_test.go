package clog_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusaconfig"
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

type testConfig struct {
	level slog.Level
	b     *bytes.Buffer
}

func (c *testConfig) GetHTTPPort() int {
	return 8080
}
func (c *testConfig) GetBaseURL() string {
	return "https://example.com"
}
func (c *testConfig) LoggerConfig() clog.Config {
	return &testLoggerConfig{level: c.level, b: c.b}
}

type testLoggerConfig struct {
	level slog.Level
	b     *bytes.Buffer
}

func (c *testLoggerConfig) Handler() (slog.Handler, error) {
	return slog.NewTextHandler(c.b, &slog.HandlerOptions{
		Level: c.level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Time(slog.TimeKey, time.Time{})
			}
			return a
		},
	}), nil
}

func TestNewDefaultConfig(t *testing.T) {
	c := clog.NewDefaultConfig(slog.LevelWarn)
	assert.NotNil(t, c)
	h, err := c.Handler()
	assert.NoError(t, err)
	assert.NotNil(t, h)
}

func TestDefaultHandler(t *testing.T) {
	h := clog.DefaultHandler(slog.LevelDebug)
	assert.NotNil(t, h)
}

func TestRegister(t *testing.T) {
	t.Run("with logger config", func(t *testing.T) {
		b := bytes.NewBuffer(nil)
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() salusaconfig.Config {
			return &testConfig{level: slog.LevelInfo, b: b}
		})

		err := clog.Register(ctx)
		assert.NoError(t, err)

		logger, err := di.Resolve[*clog.RootLogger](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, logger)

		(*slog.Logger)(logger).Warn("test")
		assert.Contains(t, b.String(), "level=WARN msg=test")
	})

	t.Run("without logger config", func(t *testing.T) {
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		di.RegisterSingleton(ctx, func() salusaconfig.Config {
			return plainConfig{}
		})

		err := clog.Register(ctx)
		assert.NoError(t, err)

		logger, err := di.Resolve[*clog.RootLogger](ctx)
		assert.NoError(t, err)
		assert.NotNil(t, logger)
	})
}

type plainConfig struct{}

func (plainConfig) GetHTTPPort() int {
	return 8080
}
func (plainConfig) GetBaseURL() string {
	return "https://example.com"
}

func TestUse(t *testing.T) {
	t.Run("resolves", func(t *testing.T) {
		ctx, _ := register()
		logger := clog.Use(ctx)
		assert.NotNil(t, logger)
	})

	t.Run("fallback to default", func(t *testing.T) {
		ctx := di.ContextWithDependencyProvider(
			context.Background(),
			di.NewDependencyProvider(),
		)
		logger := clog.Use(ctx)
		assert.NotNil(t, logger)
	})
}
