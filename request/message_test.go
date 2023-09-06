package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	msg := getMessage("max", &MessageOptions{
		Attribute: "foo",
		Value:     1,
		Arguments: []string{"5"},
	})

	assert.Equal(t, msg, "The foo must not be greater than 5.")
}
