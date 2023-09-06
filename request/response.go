package request

import (
	"io"
	"net/http"
)

type Responder interface {
	Respond(w http.ResponseWriter) error
}

type Response struct {
	body    io.Reader
	status  int
	headers map[string]string
}

var _ Responder = &Response{}

func NewResponse(body io.Reader) *Response {
	return &Response{
		body:    body,
		headers: map[string]string{},
	}
}

func (r *Response) SetStatus(status int) *Response {
	r.status = status
	return r
}

func (r *Response) AddHeader(key, value string) *Response {
	r.headers[key] = value
	return r
}

func (r *Response) Respond(w http.ResponseWriter) error {
	if r.status != 0 {
		w.WriteHeader(r.status)
	}
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}

	_, err := io.Copy(w, r.body)
	return err
}
