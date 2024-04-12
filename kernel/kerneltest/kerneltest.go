package kerneltest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
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

func (k *TestKernel[T]) Get(target string) *HttpResult {
	return k.HandleRequest(k.NewRequest(http.MethodGet, k.url(target), http.NoBody))
}
func (k *TestKernel[T]) GetJSON(target string) *HttpResult {
	r := k.NewRequest(http.MethodGet, k.url(target), http.NoBody)
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	return k.HandleRequest(r)
}

func (k *TestKernel[T]) Post(target string, body io.Reader) *HttpResult {
	return k.HandleRequest(k.NewRequest(http.MethodPost, k.url(target), body))
}

func (k *TestKernel[T]) url(target string) string {
	b := k.kernel.Config().GetBaseURL()
	return strings.TrimSuffix(b, "/") + "/" + strings.TrimPrefix(target, "/")
}

func (k *TestKernel[T]) NewRequest(method, target string, body io.Reader) *http.Request {
	return httptest.NewRequest(method, target, body).WithContext(k.ctx)
}
func (k *TestKernel[T]) PostJSON(target string, body any) *HttpResult {
	b := &bytes.Buffer{}
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		panic("body json: " + err.Error())
	}
	r := k.NewRequest(http.MethodPost, target, b)
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	return k.HandleRequest(r)
}

func (k *TestKernel[T]) HandleRequest(r *http.Request) *HttpResult {
	h := k.kernel.RootHandler(k.ctx)

	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	return &HttpResult{
		response: w.Result(),
		t:        k.t,
	}
}
