package builder_test

import (
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Clone(t *testing.T) {
	b := builder.NewBuilder().Where("a", "!=", 5)

	assert.Equal(t, b, b.Clone())
}
