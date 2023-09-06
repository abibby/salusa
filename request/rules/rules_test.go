package rules

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Rule string
	*ValidationOptions
	Valid bool
}

func runTests(t *testing.T, cases map[string]TestCase) {
	for name, d := range cases {
		t.Run(name, func(t *testing.T) {
			rule, ok := GetRule(d.Rule)

			if assert.True(t, ok, "could not find rule") {
				valid := rule(d.ValidationOptions)

				assert.Equal(t, d.Valid, valid)
			}
		})
	}
}
