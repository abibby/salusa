package integration_test

import (
	"testing"

	"github.com/abibby/salusa/static/template/test"
)

func TestIntegration(t *testing.T) {
	test.Kernel(t).
		GetJSON("/user").
		AssertStatus(200).
		AssertJSONString(`{
			"users": []
		}`)
}
