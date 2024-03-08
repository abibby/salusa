package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/static/template/app"
	"github.com/abibby/salusa/static/template/app/models/factories"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	kernel.RunTest(t, app.Kernel, "", func(t *testing.T, h http.Handler, db *sqlx.DB) {
		u := factories.UserFactory.Create(db)
		req := httptest.NewRequest("GET", fmt.Sprintf("/user?user_id=%d", u.ID), http.NoBody)
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)

		assert.Equal(t, 200, resp.Result().StatusCode)
		// assert.Equal(t, fmt.Sprintf(`{"user":{"id":%d,"username":"%s"}}`+"\n", u.ID, u.Username), resp.Body.String())
		assertJsonEqual(t, map[string]any{"user": u}, resp)
	})
}

func assertJsonEqual(t *testing.T, expected any, resp *httptest.ResponseRecorder) {
	b, err := json.Marshal(expected)
	if err != nil {
		t.Errorf("json marshal failed: %v", err)
	}
	var expectedMap any
	var actualMap any
	json.Unmarshal(b, &expectedMap)
	json.Unmarshal(resp.Body.Bytes(), &actualMap)
	assert.Equal(t, expectedMap, actualMap)
}
