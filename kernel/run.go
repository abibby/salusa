package kernel

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/abibby/salusa/clog"
)

func (k *Kernel) Run(ctx context.Context) error {

	go k.RunServices(ctx)

	return k.RunHttpServer(ctx)
}

func (k *Kernel) RunHttpServer(ctx context.Context) error {
	slog.Info(fmt.Sprintf("http://localhost:%d", k.cfg.GetHTTPPort()))

	handler := k.rootHandler(ctx)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", k.cfg.GetHTTPPort()),
		Handler: handler,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}
	return server.ListenAndServe()
}

func (k *Kernel) RunServices(ctx context.Context) {
	for _, s := range k.services {
		ctx := clog.With(ctx, slog.String("service", s.Name()))
		go func(ctx context.Context, s Service) {
			for {
				err := k.dp.Fill(ctx, s)
				if err != nil {
					k.Logger(ctx).Error("service dependency injection failed", slog.Any("error", err))
				}
				err = s.Run(ctx)
				if err != nil {
					k.Logger(ctx).Error("service failed", slog.Any("error", err))
				}
				if _, ok := s.(Restarter); !ok {
					return
				}
			}
		}(ctx, s)
	}
}
