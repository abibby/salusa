package request

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseBuilder(t *testing.T) {
	r := NewResponse(bytes.NewBufferString("hello"))
	assert.Equal(t, http.Header{}, r.Headers())
	assert.Equal(t, 0, r.Status())

	r.SetStatus(http.StatusTeapot)
	assert.Equal(t, http.StatusTeapot, r.Status())

	r.AddHeader("X-Test", "value")
	assert.Equal(t, "value", r.Header.Get("X-Test"))

	resp := r.Build()
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)

	rw := httptest.NewRecorder()
	err := r.Respond(rw, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, rw.Code)
	assert.Equal(t, "hello", rw.Body.String())
}

func TestNewResponseWrapsReadCloser(t *testing.T) {
	r := NewResponse(io.NopCloser(bytes.NewBufferString("x")))
	assert.Equal(t, "x", func() string {
		b, _ := io.ReadAll(r.Body)
		return string(b)
	}())
}

func TestJSONResponse(t *testing.T) {
	r := NewJSONResponse(map[string]any{"foo": "bar"})
	r.SetStatus(http.StatusCreated)
	r.AddHeader("X-Test", "value")

	rw := httptest.NewRecorder()
	err := r.Respond(rw, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rw.Code)
	assert.Equal(t, "application/json", rw.Header().Get("Content-Type"))
	assert.Equal(t, "value", rw.Header().Get("X-Test"))
	assert.JSONEq(t, `{"foo": "bar"}`, rw.Body.String())
}

func TestNewHTMLResponse(t *testing.T) {
	r := NewHTMLResponse([]byte("<h1>hello</h1>"))

	rw := httptest.NewRecorder()
	err := r.Respond(rw, nil)
	assert.NoError(t, err)
	assert.Equal(t, "text/html", rw.Header().Get("Content-Type"))
	assert.Equal(t, "<h1>hello</h1>", rw.Body.String())
}

func TestGetResponder(t *testing.T) {
	err := NewHTTPError(nil, http.StatusTeapot)
	responder, ok := getResponder(err)
	assert.True(t, ok)
	assert.Same(t, err, responder)

	_, ok = getResponder(nil)
	assert.False(t, ok)
}
