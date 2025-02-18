package handlertest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
)

type RequestBuilder struct {
	t       *testing.T
	ctx     context.Context
	handler http.Handler
	header  http.Header
}

func New(ctx context.Context, t *testing.T, h http.Handler) *RequestBuilder {
	return &RequestBuilder{
		t:       t,
		ctx:     ctx,
		handler: h,
		header:  http.Header{},
	}
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	rb.header.Add(key, value)
	return rb
}
func (rb *RequestBuilder) WithJSONHeaders() *RequestBuilder {
	return rb.WithHeader("Accept", "application/json").
		WithHeader("Content-Type", "application/json")
}
func (rb *RequestBuilder) build(method, target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body).WithContext(rb.ctx)
	req.Header = rb.header.Clone()
	return req
}

func jsonReader(body any) io.Reader {
	r, w := io.Pipe()
	go func() {
		err := json.NewEncoder(w).Encode(body)
		r.CloseWithError(err)
	}()
	return r
}
func (rb *RequestBuilder) handle(r *http.Request) *HttpResult {

	w := httptest.NewRecorder()

	rb.handler.ServeHTTP(w, r)

	return &HttpResult{
		response: w.Result(),
		t:        rb.t,
	}
}

func (rb *RequestBuilder) url(target string) string {
	return "https://" + path.Join("example.com", target)
}
func (rb *RequestBuilder) Get(target string) *HttpResult {
	return rb.handle(rb.build(http.MethodGet, rb.url(target), http.NoBody))
}
func (rb *RequestBuilder) GetJSON(target string) *HttpResult {
	return rb.WithJSONHeaders().Get(target)
}

func (rb *RequestBuilder) Post(target string, body io.Reader) *HttpResult {
	return rb.handle(rb.build(http.MethodPost, rb.url(target), body))
}
func (rb *RequestBuilder) PostJSON(target string, body any) *HttpResult {
	return rb.WithJSONHeaders().Post(target, jsonReader(body))
}

func (rb *RequestBuilder) Put(target string, body io.Reader) *HttpResult {
	return rb.handle(rb.build(http.MethodPut, rb.url(target), body))
}
func (rb *RequestBuilder) PutJSON(target string, body any) *HttpResult {
	return rb.WithJSONHeaders().Put(target, jsonReader(body))
}

func (rb *RequestBuilder) Delete(target string, body io.Reader) *HttpResult {
	return rb.handle(rb.build(http.MethodDelete, rb.url(target), body))
}
func (rb *RequestBuilder) DeleteJSON(target string, body any) *HttpResult {
	return rb.WithJSONHeaders().Delete(target, jsonReader(body))
}

func (rb *RequestBuilder) Patch(target string, body io.Reader) *HttpResult {
	return rb.handle(rb.build(http.MethodPatch, rb.url(target), body))
}
func (rb *RequestBuilder) PatchJSON(target string, body any) *HttpResult {
	return rb.WithJSONHeaders().Patch(target, jsonReader(body))
}
