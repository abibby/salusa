package integration_test

import (
	"testing"

	"abibby.com/salusa/static/template/test"
)

func TestIntegration(t *testing.T) {
	test.Kernel(t).
		GetJSON("/api/user").
		AssertStatus2XX().
		AssertJSONString(`{
			"users": []
		}`)
}
