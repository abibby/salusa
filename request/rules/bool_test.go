package rules

import (
	"testing"
)

func TestBool(t *testing.T) {
	data := map[string]TestCase{
		"accepted-pass": {"accepted", &ValidationOptions{Value: true}, true},
		"accepted-fail": {"accepted", &ValidationOptions{Value: false}, false},
		"declined-pass": {"declined", &ValidationOptions{Value: false}, true},
		"declined-fail": {"declined", &ValidationOptions{Value: true}, false},
	}

	runTests(t, data)
}
