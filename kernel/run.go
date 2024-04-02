package kernel

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/di"
)

func (k *Kernel) Run(ctx context.Context) error {

	go k.RunServices(ctx)

	return k.RunHttpServer(ctx)
}

func (k *Kernel) RunHttpServer(ctx context.Context) error {
	k.Logger(ctx).Info(fmt.Sprintf("listening at http://localhost:%d", k.cfg.GetHTTPPort()))

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
				err := di.Fill(ctx, s)
				if err != nil {
					k.Logger(ctx).Error("service dependency injection failed", slog.Any("error", err))
					return
				}
				err = s.Run(ctx)
				if err == nil {
					return
				}
				k.Logger(ctx).Error("service failed", slog.Any("error", err))
				r, ok := s.(Restarter)
				if !ok {
					return
				}
				r.Restart()
			}
		}(ctx, s)
	}
}
