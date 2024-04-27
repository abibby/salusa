package kerneltest

import (
	"encoding/json"
	"fmt"
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
func (r *HttpResult) AssertStatusRange(min, max int) *HttpResult {
	msg := fmt.Sprintf("Statuses must be between %d and %d", min, max)
	assert.GreaterOrEqual(r.t, r.response.StatusCode, min, msg)
	assert.LessOrEqual(r.t, r.response.StatusCode, max, msg)
	return r
}
func (r *HttpResult) AssertStatusOK() *HttpResult {
	return r.AssertStatusRange(200, 399)
}
func (r *HttpResult) AssertStatus2XX() *HttpResult {
	return r.AssertStatusRange(200, 299)
}
func (r *HttpResult) AssertStatus3XX() *HttpResult {
	return r.AssertStatusRange(300, 399)
}
func (r *HttpResult) AssertStatus4XX() *HttpResult {
	return r.AssertStatusRange(400, 499)
}
func (r *HttpResult) AssertStatus5XX() *HttpResult {
	return r.AssertStatusRange(500, 599)
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
