package jobs

import (
	"context"
	"log/slog"
	"testing"

	"abibby.com/salusa/di"
	"abibby.com/salusa/static/template/app/events"
	"github.com/stretchr/testify/require"
)

func TestLogJob_Handle(t *testing.T) {
	ctx := di.TestDependencyProviderContext()
	di.RegisterSingleton(ctx, func() *slog.Logger {
		return slog.Default()
	})

	l := &LogJob{}
	require.NoError(t, di.Fill(ctx, l))

	err := l.Handle(context.Background(), &events.LogEvent{Message: "test"})
	require.NoError(t, err)
}
