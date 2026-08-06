package dialects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultPanics(t *testing.T) {
	old := defaultDialect
	defaultDialect = func() Dialect { panic("no default dialect set") }
	defer func() {
		defaultDialect = old
	}()
	assert.Panics(t, func() {
		New()
	})
}
