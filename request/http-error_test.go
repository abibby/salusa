package request

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusError(t *testing.T) {
	err := ErrStatusNotFound
	assert.Equal(t, "http 404: Not Found", err.Error())

	rw := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	respondErr := err.Respond(rw, req)
	assert.NoError(t, respondErr)
	assert.Equal(t, http.StatusNotFound, rw.Code)
	assert.Contains(t, rw.Body.String(), "Not Found")
}

func TestHTTPError(t *testing.T) {
	inner := errors.New("boom")
	e := NewHTTPError(inner, http.StatusTeapot)

	assert.Equal(t, "http 418: boom", e.Error())
	assert.Same(t, inner, e.Unwrap())
	assert.Equal(t, http.StatusTeapot, e.Status())

	e.WithStack()
	assert.NotNil(t, e.stack)
}

func TestHTTPErrorRespondJSON(t *testing.T) {
	t.Run("without stack", func(t *testing.T) {
		e := NewHTTPError(errors.New("boom"), http.StatusBadRequest)
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		err := e.Respond(rw, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "boom")
	})

	t.Run("with stack on 500", func(t *testing.T) {
		e := NewHTTPError(errors.New("boom"), http.StatusInternalServerError)
		e.WithStack()
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		err := e.Respond(rw, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "go_routine")
	})

	t.Run("with validation error fields", func(t *testing.T) {
		verr := ValidationError{"foo": []string{"is required"}}
		e := NewHTTPError(verr, http.StatusUnprocessableEntity)
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		err := e.Respond(rw, req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
		assert.Contains(t, rw.Body.String(), "is required")
	})
}

func TestHTTPErrorRespondHTML(t *testing.T) {
	t.Run("plain error", func(t *testing.T) {
		e := NewHTTPError(errors.New("boom"), http.StatusBadRequest)
		rw := httptest.NewRecorder()
		err := e.Respond(rw, httptest.NewRequest("GET", "/", nil))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
		assert.Contains(t, rw.Body.String(), "Error: 400 Bad Request")
		assert.Contains(t, rw.Body.String(), "<h2>boom</h2>")
	})

	t.Run("html error", func(t *testing.T) {
		e := NewHTTPError(ValidationError{"foo": []string{"bad"}}, http.StatusBadRequest)
		rw := httptest.NewRecorder()
		err := e.Respond(rw, httptest.NewRequest("GET", "/", nil))
		assert.NoError(t, err)
		assert.Contains(t, rw.Body.String(), "Validation Error")
	})

	t.Run("with stack", func(t *testing.T) {
		e := NewHTTPError(errors.New("boom"), http.StatusInternalServerError)
		e.WithStack()
		rw := httptest.NewRecorder()
		err := e.Respond(rw, httptest.NewRequest("GET", "/", nil))
		assert.NoError(t, err)
		assert.Contains(t, rw.Body.String(), "Error: 500 Internal Server Error")
	})

	t.Run("with stack frames", func(t *testing.T) {
		dir := t.TempDir()
		file := filepath.Join(dir, "main.go")
		var content strings.Builder
		for i := 0; i < 20; i++ {
			fmt.Fprintf(&content, "line %d\n", i)
		}
		assert.NoError(t, os.WriteFile(file, []byte(content.String()), 0o644))

		e := NewHTTPError(errors.New("boom"), http.StatusInternalServerError)
		e.stack = []byte("goroutine 1 [running]:\n" +
			"panic({0xdeadbeef})\n" +
			"\t/usr/local/go/src/runtime/panic.go:884 +0x212\n" +
			"main.foo()\n" +
			"\t" + file + ":10 +0x20\n" +
			"main.bar()\n" +
			"\t" + filepath.Join(dir, "missing.go") + ":5 +0x10\n")
		rw := httptest.NewRecorder()
		err := e.Respond(rw, httptest.NewRequest("GET", "/", nil))
		assert.NoError(t, err)
		assert.Contains(t, rw.Body.String(), "main.foo()")
		assert.Contains(t, rw.Body.String(), "line 10")
		assert.Contains(t, rw.Body.String(), "Error: open")
	})
}

func TestParseStack(t *testing.T) {
	stack := []byte("goroutine 1 [running]:\n" +
		"panic({0xdeadbeef})\n" +
		"\t/usr/local/go/src/runtime/panic.go:884 +0x212\n" +
		"main.foo()\n" +
		"\t/home/user/main.go:10 +0x20\n" +
		"main.main()\n" +
		"\t/home/user/main.go:15 +0x16\n")

	parsed := parseStack(stack)
	assert.Equal(t, "goroutine 1 [running]:", parsed.GoRoutine)
	assert.Len(t, parsed.Stack, 2)
	assert.Equal(t, "main.foo()", parsed.Stack[0].Call)
	assert.Equal(t, "/home/user/main.go", parsed.Stack[0].File)
	assert.Equal(t, 10, parsed.Stack[0].Line)
	assert.Equal(t, 0x20, parsed.Stack[0].Extra)

	t.Run("no extra part", func(t *testing.T) {
		stack := []byte("goroutine 1 [running]:\n" +
			"panic({0xdeadbeef})\n" +
			"\t/usr/local/go/src/runtime/panic.go:884\n" +
			"main.foo()\n" +
			"\t/home/user/main.go:10\n")
		parsed := parseStack(stack)
		assert.Len(t, parsed.Stack, 1)
		assert.Equal(t, -1, parsed.Stack[0].Extra)
		assert.Equal(t, 10, parsed.Stack[0].Line)
	})

	t.Run("bad line number", func(t *testing.T) {
		stack := []byte("goroutine 1 [running]:\n" +
			"panic({0xdeadbeef})\n" +
			"\t/usr/local/go/src/runtime/panic.go:884 +0x212\n" +
			"main.foo()\n" +
			"\t/home/user/main.go:not-a-number +0x20\n")
		parsed := parseStack(stack)
		assert.Len(t, parsed.Stack, 1)
		assert.Equal(t, -1, parsed.Stack[0].Line)
	})
}

func TestParsInt(t *testing.T) {
	v, err := parsInt("+0x212")
	assert.NoError(t, err)
	assert.Equal(t, 0x212, v)

	v, err = parsInt("-5")
	assert.NoError(t, err)
	assert.Equal(t, -5, v)

	v, err = parsInt("123")
	assert.NoError(t, err)
	assert.Equal(t, 123, v)

	_, err = parsInt("abc")
	assert.Error(t, err)
}

func TestErrorHandler(t *testing.T) {
	t.Run("plain error becomes 500", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		ErrorHandler(errors.New("boom")).ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "boom")
	})

	t.Run("responder error", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		ErrorHandler(ErrStatusTeapot).ServeHTTP(rw, req)
		assert.Equal(t, http.StatusTeapot, rw.Code)
	})

	t.Run("failing responder", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		ErrorHandler(failingResponderErr{err: errors.New("inner")}).ServeHTTP(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
	})
}

type failingResponderErr struct {
	err error
}

func (e failingResponderErr) Error() string { return "failing responder" }
func (e failingResponderErr) Respond(w http.ResponseWriter, r *http.Request) error {
	return e.err
}

func TestHTTPErrorRespondClosesPipe(t *testing.T) {
	e := NewHTTPError(errors.New("boom"), http.StatusBadRequest)
	rw := httptest.NewRecorder()
	err := e.Respond(rw, httptest.NewRequest("GET", "/", nil))
	assert.NoError(t, err)
	assert.True(t, strings.Contains(rw.Body.String(), "boom"))
}
