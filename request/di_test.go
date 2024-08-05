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
		Request *http.Request `inject:""`
	}
	ctx := di.TestDependencyProviderContext()
	err := request.Register(ctx)
	assert.NoError(t, err)

	httpRequest := httptest.
		NewRequest("GET", "http://0.0.0.0/", http.NoBody).
		WithContext(ctx)

	h := request.Handler(func(r *Request) (any, error) {
		assert.Same(t, httpRequest, r.Request)
		return nil, nil
	})

	h.ServeHTTP(
		httptest.NewRecorder(),
		httpRequest,
	)
}

func TestInjectResponseWriter(t *testing.T) {
	type Request struct {
		ResponseWriter http.ResponseWriter `inject:""`
	}

	ctx := di.TestDependencyProviderContext()
	err := request.Register(ctx)
	assert.NoError(t, err)

	rw := httptest.NewRecorder()

	h := request.Handler(func(r *Request) (any, error) {
		assert.Same(t, rw, r.ResponseWriter)
		return nil, nil
	})

	h.ServeHTTP(
		rw,
		httptest.
			NewRequest("GET", "http://0.0.0.0/", http.NoBody).
			WithContext(ctx),
	)
}
