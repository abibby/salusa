package kerneltest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/abibby/salusa/kernel"
)

type testKernel interface {
	HandleRequest(r *http.Request) *HttpResult
	url(target string) string
}

type RequestBuilder struct {
	kernel testKernel
	ctx    context.Context
	header http.Header
}

func NewRequestBuilder[T kernel.KernelConfig](k *TestKernel[T]) *RequestBuilder {
	return &RequestBuilder{
		kernel: k,
		ctx:    k.ctx,
		header: http.Header{},
	}
}

func (rb *RequestBuilder) WithHeader(key, value string) *RequestBuilder {
	rb.header.Add(key, value)
	return rb
}
func (rb *RequestBuilder) NewRequest(method, target string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method, target, body).WithContext(rb.ctx)
	req.Header = rb.header
	return req
}
func (rb *RequestBuilder) Get(target string) *HttpResult {
	return rb.kernel.HandleRequest(rb.NewRequest(http.MethodGet, rb.kernel.url(target), http.NoBody))
}
func (rb *RequestBuilder) GetJSON(target string) *HttpResult {
	return rb.WithHeader("Accept", "application/json").
		WithHeader("Content-Type", "application/json").
		Get(target)
}

func (rb *RequestBuilder) Post(target string, body io.Reader) *HttpResult {
	return rb.kernel.HandleRequest(rb.NewRequest(http.MethodPost, rb.kernel.url(target), body))
}

func (rb *RequestBuilder) PostJSON(target string, body any) *HttpResult {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		panic("body json: " + err.Error())
	}

	return rb.WithHeader("Accept", "application/json").
		WithHeader("Content-Type", "application/json").
		Post(target, b)
}
