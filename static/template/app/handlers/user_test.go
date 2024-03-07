package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/static/template/app"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	kernel.RunTest(t, app.Kernel, "", func(t *testing.T, h http.Handler) {
		req := httptest.NewRequest("GET", "/user?user_id=2", http.NoBody)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Result().StatusCode)
		assert.Equal(t, "", resp.Body.String())
	})
}
