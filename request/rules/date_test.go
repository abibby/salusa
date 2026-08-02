package rules

import (
	"testing"
	"time"
)

var (
	jan1 = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	jan2 = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
)

func TestDate(t *testing.T) {
	data := map[string]TestCase{
		"after-pass":             {"after", &ValidationOptions{Value: jan2, Arguments: []string{"2020-01-01T00:00:00Z"}}, true},
		"after-fail":             {"after", &ValidationOptions{Value: jan1, Arguments: []string{"2020-01-01T00:00:00Z"}}, false},
		"after_or_equal-pass":    {"after_or_equal", &ValidationOptions{Value: jan1, Arguments: []string{"2020-01-01T00:00:00Z"}}, true},
		"after_or_equal-fail":    {"after_or_equal", &ValidationOptions{Value: jan1, Arguments: []string{"2020-01-02T00:00:00Z"}}, false},
		"before-pass":            {"before", &ValidationOptions{Value: jan1, Arguments: []string{"2020-01-02T00:00:00Z"}}, true},
		"before-fail":            {"before", &ValidationOptions{Value: jan2, Arguments: []string{"2020-01-01T00:00:00Z"}}, false},
		"before_or_equal-pass":   {"before_or_equal", &ValidationOptions{Value: jan2, Arguments: []string{"2020-01-02T00:00:00Z"}}, true},
		"before_or_equal-fail":   {"before_or_equal", &ValidationOptions{Value: jan2, Arguments: []string{"2020-01-01T00:00:00Z"}}, false},
		"time-pointer-pass":      {"after", &ValidationOptions{Value: ptr(jan2), Arguments: []string{"2020-01-01T00:00:00Z"}}, true},
		"time-invalid-arg-pass":  {"after", &ValidationOptions{Value: jan2, Arguments: []string{"not a time"}}, true},
	}

	runTests(t, data)
}
