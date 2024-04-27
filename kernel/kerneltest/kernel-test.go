package kerneltest

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type TestKernel[T kernel.KernelConfig] struct {
	kernel *kernel.Kernel
	config T
	ctx    context.Context
	t      *testing.T
}

func NewTestKernelFactory[T kernel.KernelConfig](k *kernel.Kernel, cfg T) func(t *testing.T) *TestKernel[T] {
	return func(t *testing.T) *TestKernel[T] {
		ctx := di.TestDependencyProviderContext()
		k := kernel.Config(func() T { return cfg })(k)
		err := k.Bootstrap(ctx)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
		return &TestKernel[T]{
			kernel: k,
			config: cfg,
			ctx:    ctx,
			t:      t,
		}
	}
}

func (k *TestKernel[T]) url(target string) string {
	b := k.kernel.Config().GetBaseURL()
	return strings.TrimSuffix(b, "/") + "/" + strings.TrimPrefix(target, "/")
}

func (k *TestKernel[T]) WithHeader(key, value string) *RequestBuilder {
	return NewRequestBuilder(k).WithHeader(key, value)
}
func (k *TestKernel[T]) WithJSONHeaders() *RequestBuilder {
	return NewRequestBuilder(k).WithJSONHeaders()
}

func (k *TestKernel[T]) Get(target string) *HttpResult {
	return NewRequestBuilder(k).Get(target)
}
func (k *TestKernel[T]) GetJSON(target string) *HttpResult {
	return NewRequestBuilder(k).GetJSON(target)
}
func (k *TestKernel[T]) Post(target string, body io.Reader) *HttpResult {
	return NewRequestBuilder(k).Post(target, body)
}
func (k *TestKernel[T]) PostJSON(target string, body any) *HttpResult {
	return NewRequestBuilder(k).PostJSON(target, body)
}
func (k *TestKernel[T]) Put(target string, body io.Reader) *HttpResult {
	return NewRequestBuilder(k).Put(target, body)
}
func (k *TestKernel[T]) PutJSON(target string, body any) *HttpResult {
	return NewRequestBuilder(k).PutJSON(target, body)
}
func (k *TestKernel[T]) Patch(target string, body io.Reader) *HttpResult {
	return NewRequestBuilder(k).Patch(target, body)
}
func (k *TestKernel[T]) PatchJSON(target string, body any) *HttpResult {
	return NewRequestBuilder(k).PatchJSON(target, body)
}
func (k *TestKernel[T]) Delete(target string, body io.Reader) *HttpResult {
	return NewRequestBuilder(k).Delete(target, body)
}
func (k *TestKernel[T]) DeleteJSON(target string, body any) *HttpResult {
	return NewRequestBuilder(k).DeleteJSON(target, body)
}

func (k *TestKernel[T]) HandleRequest(r *http.Request) *HttpResult {
	h := k.kernel.RootHandler(k.ctx)
	spew.Dump(di.GetDependencyProvider(k.ctx))
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	return &HttpResult{
		response: w.Result(),
		t:        k.t,
	}
}
