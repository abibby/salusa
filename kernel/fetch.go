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

type ResponseWriter struct {
	header http.Header
	writer io.Writer
	status int
}

func NewStdoutResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header: http.Header{},
		writer: os.Stdout,
		status: 0,
	}
}

func (w *ResponseWriter) Ok() bool {
	return w.status >= 200 && w.status < 300
}

func (w *ResponseWriter) Status() int {
	return w.status
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

// Header implements http.ResponseWriter.
func (s *ResponseWriter) Header() http.Header {
	return s.header
}

// Write implements http.ResponseWriter.
func (s *ResponseWriter) Write(b []byte) (int, error) {
	if s.status == 0 {
		s.WriteHeader(200)
	}
	return s.writer.Write(b)
}

// WriteHeader implements http.ResponseWriter.
func (s *ResponseWriter) WriteHeader(statusCode int) {
	if s.status != 0 {
		return
	}
	for k, vs := range s.header {
		for _, v := range vs {
			fmt.Fprintf(s.writer, "< header: %s: %s\n", k, v)
		}
	}
	fmt.Fprintf(s.writer, "< status: %d\n", statusCode)
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
