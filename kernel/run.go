package kernel

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strings"

	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/di"
	"github.com/spf13/pflag"
)

type StdoutResponseWriter struct {
	header http.Header
	status int
}

func NewStdoutResponseWriter() *StdoutResponseWriter {
	return &StdoutResponseWriter{
		header: http.Header{},
		status: 200,
	}
}

func (w *StdoutResponseWriter) Ok() bool {
	return w.status >= 200 && w.status < 300
}

func (w *StdoutResponseWriter) Status() int {
	return w.status
}

var _ http.ResponseWriter = (*StdoutResponseWriter)(nil)

// Header implements http.ResponseWriter.
func (s *StdoutResponseWriter) Header() http.Header {
	return s.header
}

// Write implements http.ResponseWriter.
func (s *StdoutResponseWriter) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

// WriteHeader implements http.ResponseWriter.
func (s *StdoutResponseWriter) WriteHeader(statusCode int) {
	s.status = statusCode
}

func (k *Kernel) Run(ctx context.Context) error {
	validate := pflag.BoolP("validate", "v", false, "validate di")
	fetch := pflag.String("fetch", "", "run a single request and print the result to stdout")
	method := pflag.StringP("method", "m", "get", "method for fetch")
	headers := pflag.StringArrayP("header", "h", []string{}, "header")

	pflag.Parse()

	err := k.Validate(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	if *validate {
		if err != nil {
			os.Exit(1)
		}
		return nil
	}

	if *fetch != "" {
		return k.runFetch(ctx, *fetch, *method, *headers)
	}

	go k.RunServices(ctx)

	return k.RunHttpServer(ctx)
}
func (k *Kernel) runFetch(ctx context.Context, uri, method string, headers []string) error {
	var u *url.URL
	if strings.HasPrefix(uri, "/") {
		uri = "http://localhost" + uri
	} else if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		// noop
	} else {
		uri = "http://localhost/" + uri
	}

	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	h := k.handlerWithMiddleware()
	r, err := http.NewRequest(strings.ToUpper(method), u.String(), http.NoBody)
	if err != nil {
		return err
	}
	u.Host = ""
	u.Scheme = ""
	r.RequestURI = u.String()
	r = r.WithContext(ctx)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")
	for _, header := range headers {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Invalid header %s", header)
		}
		r.Header.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	w := NewStdoutResponseWriter()
	h.ServeHTTP(w, r)
	if !w.Ok() {
		os.Exit(w.status / 100)
	}
	return nil
}
func (k *Kernel) handlerWithMiddleware() http.Handler {
	h := k.rootHandler

	for _, m := range k.globalMiddleware {
		h = m(h)
	}
	return h
}

func (k *Kernel) RunHttpServer(ctx context.Context) error {
	clog.Use(ctx).Info(fmt.Sprintf("listening at http://localhost:%d", k.cfg.GetHTTPPort()))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", k.cfg.GetHTTPPort()),
		Handler: k.handlerWithMiddleware(),
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	k.addCloser(server)

	return server.ListenAndServe()
}

func (k *Kernel) RunServices(ctx context.Context) {
	for _, s := range k.services {
		ctx := clog.With(ctx, slog.String("service", s.Name()))
		go func(ctx context.Context, s Service) {
			for {
				err := di.Fill(ctx, s)
				if err != nil {
					clog.Use(ctx).Error("service dependency injection failed", slog.Any("error", err))
					return
				}
				err = s.Run(ctx)
				if err == nil {
					return
				}
				clog.Use(ctx).Error("service failed", slog.Any("error", err))
				r, ok := s.(Restarter)
				if !ok {
					return
				}
				r.Restart()
			}
		}(ctx, s)
	}
}

func (k *Kernel) singles(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		tries := 0
		for range c {
			if tries == 0 {
				go func() {
					log.Print("Gracefully shutting down server Ctrl+C again to force")
					for _, closer := range k.closers {
						err := closer.Close()
						if err != nil {
							slog.Error("failed to close resource",
								"err", err,
								"resource", reflect.TypeOf(closer),
							)
						}
					}
					os.Exit(1)
				}()
			} else {
				log.Print("Force shutting down server")
				os.Exit(1)
			}
			tries++
		}
	}()
}
