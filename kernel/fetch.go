package kernel

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/salusaconfig"
)

type StdoutResponseWriter struct {
	header http.Header
	status int
}

func NewStdoutResponseWriter() *StdoutResponseWriter {
	return &StdoutResponseWriter{
		header: http.Header{},
		status: 0,
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
	if s.status == 0 {
		s.WriteHeader(200)
	}
	return os.Stdout.Write(b)
}

// WriteHeader implements http.ResponseWriter.
func (s *StdoutResponseWriter) WriteHeader(statusCode int) {
	if s.status != 0 {
		return
	}
	for k, vs := range s.header {
		for _, v := range vs {
			fmt.Fprintf(os.Stdout, "< header: %s: %s\n", k, v)
		}
	}
	fmt.Fprintf(os.Stdout, "< status: %d\n", statusCode)
	s.status = statusCode
}
func newRequest(ctx context.Context, uri, method string, body string) (*http.Request, error) {
	var u *url.URL
	var err error
	if strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://") {
		u, err = url.Parse(uri)
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err := di.Resolve[salusaconfig.Config](ctx)
		if err != nil {
			return nil, err
		}
		baseURLStr := cfg.GetBaseURL()
		if baseURLStr == "" {
			baseURLStr = "http://localhost"
		} else {
			baseURLStr = strings.TrimSuffix(baseURLStr, "/")
		}

		u, err = url.Parse(baseURLStr + "/" + strings.TrimPrefix(uri, "/"))
		if err != nil {
			return nil, err
		}
	}
	host := u.Host
	u.Host = ""
	u.Scheme = ""
	requestURI := u.String()

	var bodyReader io.Reader = http.NoBody
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
	}
	r, err := http.NewRequest(strings.ToUpper(method), requestURI, bodyReader)
	if err != nil {
		return nil, err
	}
	r.Host = host
	r.RequestURI = requestURI
	return r.WithContext(ctx), nil
}
func (k *Kernel) runFetch(ctx context.Context, uri, method string, headers []string, body string) error {
	h := k.handlerWithMiddleware()

	r, err := newRequest(ctx, uri, method, body)
	if err != nil {
		return err
	}
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
