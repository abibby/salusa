package request

import (
	"errors"
	"io"
	"net/http"
)

type Responder interface {
	Respond(w http.ResponseWriter, r *http.Request) error
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

func (r *Response) Status() int {
	return r.status
}
func (r *Response) SetStatus(status int) *Response {
	r.status = status
	return r
}

func (r *Response) Headers() map[string]string {
	return r.headers
}
func (r *Response) AddHeader(key, value string) *Response {
	r.headers[key] = value
	return r
}

func (r *Response) Respond(w http.ResponseWriter, _ *http.Request) error {
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}
	if r.status != 0 {
		w.WriteHeader(r.status)
	}

	_, err := io.Copy(w, r.body)
	return err
}

func getResponder(err error) (Responder, bool) {
	var responder Responder
	var ok bool
	for err != nil {
		responder, ok = err.(Responder)
		if ok {
			return responder, true
		}
		err = errors.Unwrap(err)
	}
	return nil, false
}
