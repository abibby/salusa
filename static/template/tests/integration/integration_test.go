package integration_test

import (
	"testing"

	"github.com/abibby/salusa/static/template/test"
)

func TestIntegration(t *testing.T) {
	test.Kernel(t).
		GetJSON("/api/user").
		AssertStatusOK().
		AssertJSONString(`{
			"users": []
		}`)
}
