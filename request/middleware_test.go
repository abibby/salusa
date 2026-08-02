package request

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToError(t *testing.T) {
	assert.Nil(t, toError(nil))

	err := errors.New("boom")
	assert.Same(t, err, toError(err))
	assert.EqualError(t, toError("boom"), "boom")
	assert.EqualError(t, toError(42), "42")
}

func TestJoinedError(t *testing.T) {
	assert.Nil(t, (&errorsContainer{}).joinedError())

	one := errors.New("one")
	assert.Same(t, one, (&errorsContainer{errors: []error{one}}).joinedError())

	joined := (&errorsContainer{errors: []error{one, errors.New("two")}}).joinedError()
	assert.Error(t, joined)
}

func TestAddErrorWithoutContainer(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	assert.NotPanics(t, func() {
		addError(r, errors.New("boom"))
	})
	assert.False(t, hasHandleErrors(r))
}

func TestHandleErrors(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		handler := HandleErrors()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		rw := httptest.NewRecorder()
		handler.ServeHTTP(rw, httptest.NewRequest("GET", "/", nil))
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	t.Run("handler returns error", func(t *testing.T) {
		handler := HandleErrors()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			RespondError(w, r, errors.New("handler error"))
		}))
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		handler.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "handler error")
	})

	t.Run("panic with error", func(t *testing.T) {
		handler := HandleErrors()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("panic error"))
		}))
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		handler.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "panic error")
		assert.Contains(t, rw.Body.String(), "stack")
	})

	t.Run("panic with string", func(t *testing.T) {
		handler := HandleErrors()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("panic string")
		}))
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		handler.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "panic string")
	})

	t.Run("panic with other value", func(t *testing.T) {
		handler := HandleErrors()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(42)
		}))
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		handler.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
		assert.Contains(t, rw.Body.String(), "42")
	})

	t.Run("panic with HTTPError preserves status", func(t *testing.T) {
		handler := HandleErrors()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(NewHTTPError(errors.New("teapot error"), http.StatusTeapot))
		}))
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		handler.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusTeapot, rw.Code)
		assert.Contains(t, rw.Body.String(), "teapot error")
	})

	t.Run("custom error handler", func(t *testing.T) {
		handler := HandleErrors(func(ctx context.Context, err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte("custom: " + err.Error()))
			})
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("boom"))
		}))
		rw := httptest.NewRecorder()
		handler.ServeHTTP(rw, httptest.NewRequest("GET", "/", nil))
		assert.Equal(t, http.StatusBadGateway, rw.Code)
		assert.Contains(t, rw.Body.String(), "boom")
	})

	t.Run("custom error handler last wins", func(t *testing.T) {
		first := func(ctx context.Context, err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadGateway)
			})
		}
		second := func(ctx context.Context, err error) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			})
		}
		handler := HandleErrors(first, second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("boom"))
		}))
		rw := httptest.NewRecorder()
		handler.ServeHTTP(rw, httptest.NewRequest("GET", "/", nil))
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run("custom error handler that returns nil", func(t *testing.T) {
		handler := HandleErrors(func(ctx context.Context, err error) http.Handler {
			return nil
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(errors.New("boom"))
		}))
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Accept", "application/json")
		handler.ServeHTTP(rw, req)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}
