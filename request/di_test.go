package request_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/request"
	"github.com/stretchr/testify/assert"
)

func TestInjectRequest(t *testing.T) {
	type Request struct {
		Request *http.Request
	}
	ctx := di.TestDependencyProviderContext()
	request.Register(ctx)

	httpRequest := httptest.
		NewRequest("GET", "http://0.0.0.0/", http.NoBody).
		WithContext(ctx)

	h := request.Handler(func(r *Request) (any, error) {
		assert.Same(t, r.Request, httpRequest)
		return nil, nil
	})

	h.ServeHTTP(
		httptest.NewRecorder(),
		httpRequest,
	)
}

func TestInjectResponseWriter(t *testing.T) {
	type Request struct {
		ResponseWriter http.ResponseWriter
	}

	ctx := di.TestDependencyProviderContext()
	request.Register(ctx)

	rw := httptest.NewRecorder()

	h := request.Handler(func(r *Request) (any, error) {
		assert.Same(t, r.ResponseWriter, rw)
		return nil, nil
	})

	h.ServeHTTP(
		rw,
		httptest.
			NewRequest("GET", "http://0.0.0.0/", http.NoBody).
			WithContext(ctx),
	)
}
