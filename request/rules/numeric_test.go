package rules

import (
	"testing"
)

func TestNumeric(t *testing.T) {
	data := map[string]TestCase{
		"max-pass":         {"max", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, true},
		"max-fail":         {"max", &ValidationOptions{Value: 2, Arguments: []string{"1"}}, false},
		"min-pass":         {"min", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, true},
		"min-fail":         {"min", &ValidationOptions{Value: 0, Arguments: []string{"1"}}, false},
		"multiple_of-pass": {"multiple_of", &ValidationOptions{Value: 10, Arguments: []string{"5"}}, true},
		"multiple_of-fail": {"multiple_of", &ValidationOptions{Value: 6, Arguments: []string{"5"}}, false},
		"ptr-max-pass":     {"max", &ValidationOptions{Value: ptr(1), Arguments: []string{"1"}}, true},
		"ptr-max-fail":     {"max", &ValidationOptions{Value: ptr(2), Arguments: []string{"1"}}, false},
	}

	runTests(t, data)
}
