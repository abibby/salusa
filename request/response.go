package request

import (
	"errors"
	"io"
	"net/http"
)

type Responder interface {
	Respond(w http.ResponseWriter, r *http.Request) error
}

type ResponseBuilder struct {
	*http.Response
}

var _ Responder = &ResponseBuilder{}

func NewResponse(body io.Reader) *ResponseBuilder {
	closer, ok := body.(io.ReadCloser)
	if !ok {
		closer = io.NopCloser(body)
	}
	return &ResponseBuilder{
		Response: &http.Response{
			Header: http.Header{},
			Body:   closer,
		},
	}
}

func (r *ResponseBuilder) Status() int {
	return r.StatusCode
}
func (r *ResponseBuilder) SetStatus(status int) *ResponseBuilder {
	r.StatusCode = status
	return r
}

func (r *ResponseBuilder) Headers() http.Header {
	return r.Header
}
func (r *ResponseBuilder) AddHeader(key, value string) *ResponseBuilder {
	r.Header.Add(key, value)
	return r
}

func (r *ResponseBuilder) Respond(w http.ResponseWriter, _ *http.Request) error {
	return serveResponse(w, r.Response)
}

func getResponder(err error) (Responder, bool) {
	var responder interface {
		error
		Responder
	}
	if errors.As(err, &responder) {
		return responder, true
	}
	return nil, false
}
