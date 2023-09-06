package request

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInjectRequest(t *testing.T) {
	type Request struct {
		Request *http.Request
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/", http.NoBody)

	h := Handler(func(r *Request) (any, error) {
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

	rw := httptest.NewRecorder()

	h := Handler(func(r *Request) (any, error) {
		assert.Same(t, r.ResponseWriter, rw)
		return nil, nil
	})

	h.ServeHTTP(
		rw,
		httptest.NewRequest("GET", "http://0.0.0.0/", http.NoBody),
	)
}
