package kernel

import (
	"context"
	"log/slog"

	"github.com/abibby/salusa/di"
)

func (k *Kernel) Logger(ctx context.Context) *slog.Logger {
	logger, err := di.Resolve[*slog.Logger](ctx, k.dp)
	if err != nil {
		logger = slog.Default()
		logger.Warn("no logger in di", slog.Any("error", err))
	}
	return logger
}
