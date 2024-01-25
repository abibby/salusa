package request_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abibby/salusa/request"
	"github.com/stretchr/testify/assert"
)

func TestInjectRequest(t *testing.T) {
	type Request struct {
		Request *http.Request
	}

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
		ResponseWriter http.ResponseWriter
	}

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

func ExampleHandler_Input() {
	type ExampleRequest struct {
		A int `query:"a"`
		B int `query:"b"`
	}
	type ExampleResponse struct {
		Sum int `json:"sum"`
	}
	h := request.Handler(func(r *ExampleRequest) (*ExampleResponse, error) {
		return &ExampleResponse{
			Sum: r.A + r.B,
		}, nil
	})

	rw := httptest.NewRecorder()

	h.ServeHTTP(
		rw,
		httptest.NewRequest("GET", "/?a=10&b=5", http.NoBody),
	)

	fmt.Println(rw.Body)
	// Output: {"sum":15}
}

func ExampleHandler_Error() {
	type ExampleRequest struct {
		A int `query:"a" validate:"min:1"`
	}
	type ExampleResponse struct {
	}

	h := request.Handler(func(r *ExampleRequest) (*ExampleResponse, error) {
		return &ExampleResponse{}, nil
	})

	rw := httptest.NewRecorder()

	h.ServeHTTP(
		rw,
		httptest.NewRequest("GET", "/?a=-1", http.NoBody),
	)

	fmt.Println(rw.Result().StatusCode)
	// Output: 422
}
