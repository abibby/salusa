package request_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/abibby/salusa/request"
)

func ExampleHandler_input() {
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
	// Output:
	// {
	//     "sum": 15
	// }
}

func ExampleHandler_error() {
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
