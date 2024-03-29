package request

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	msg, err := getMessage(context.Background(), "max", &MessageOptions{
		Attribute: "foo",
		Value:     1,
		Arguments: []string{"5"},
	})

	assert.NoError(t, err)
	assert.Equal(t, msg, "The foo must not be greater than 5.")
}
