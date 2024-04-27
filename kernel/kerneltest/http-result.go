package kerneltest

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HttpResult struct {
	t        *testing.T
	response *http.Response
	body     []byte
}

func (r *HttpResult) Body() []byte {
	if r.body == nil {
		defer r.response.Body.Close()
		b, err := io.ReadAll(r.response.Body)
		if err != nil {
			panic(err)
		}
		r.body = b
	}
	return r.body
}

func (r *HttpResult) AssertStatus(status int) *HttpResult {
	assert.Equal(r.t, status, r.response.StatusCode, "Statuses do not match")
	return r
}

func (r *HttpResult) AssertJSONString(jsonBody string) *HttpResult {
	var expected any
	err := json.Unmarshal([]byte(jsonBody), &expected)
	if err != nil {
		panic(err)
	}
	var actual any
	err = json.Unmarshal(r.Body(), &actual)
	if err != nil {
		assert.Fail(r.t, "body is not json", err.Error())
		return r
	}
	assert.Equal(r.t, expected, actual, "Statuses do not match")
	return r
}
func (r *HttpResult) AssertJSON(body any) *HttpResult {
	expectedJSON, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	return r.AssertJSONString(string(expectedJSON))
}
