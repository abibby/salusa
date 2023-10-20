package kernel

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/abibby/salusa/clog"
)

func (k *Kernel) Run(ctx context.Context) error {

	go k.RunListeners(ctx)
	go k.RunServices(ctx)

	return k.RunHttpServer(ctx)
}

func (k *Kernel) RunHttpServer(ctx context.Context) error {
	slog.Info(fmt.Sprintf("http://localhost:%d", k.port))

	handler := k.rootHandler()
	return http.ListenAndServe(fmt.Sprintf(":%d", k.port), handler)
}

func (k *Kernel) RunServices(ctx context.Context) error {
	for _, s := range k.services {
		ctx := clog.Update(ctx, func(l *slog.Logger) *slog.Logger {
			return l.With("service", s.Name())
		})
		go func(ctx context.Context, s Service) {
			for {
				err := s.Run(ctx, k)
				if err != nil {
					slog.Error("service failed", slog.Any("error", err), slog.String("service", s.Name()))
				}
				if _, ok := s.(Restarter); !ok {
					return
				}
			}
		}(ctx, s)
	}
	return nil
}
