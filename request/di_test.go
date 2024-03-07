package request_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abibby/salusa/request"
	"github.com/stretchr/testify/assert"
)

func TestInjectRequest(t *testing.T) {
	type Request struct {
		Request *http.Request `inject:""`
	}

	request.InitDI(context.Background())

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/", http.NoBody)

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
		ResponseWriter http.ResponseWriter `inject:""`
	}

	request.InitDI(context.Background())

	rw := httptest.NewRecorder()

	h := request.Handler(func(r *Request) (any, error) {
		assert.Same(t, r.ResponseWriter, rw)
		return nil, nil
	})

	h.ServeHTTP(
		rw,
		httptest.NewRequest("GET", "http://0.0.0.0/", http.NoBody),
	)
}
