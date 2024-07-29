package request

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		msg, err := getMessage(context.Background(), "max", &MessageOptions{
			Attribute: "foo",
			Value:     "",
			Arguments: []string{"5"},
		})

		assert.NoError(t, err)
		assert.Equal(t, msg, "The foo must not be greater than 5 characters.")
	})

	t.Run("number", func(t *testing.T) {
		msg, err := getMessage(context.Background(), "max", &MessageOptions{
			Attribute: "foo",
			Value:     1,
			Arguments: []string{"5"},
		})

		assert.NoError(t, err)
		assert.Equal(t, msg, "The foo must not be greater than 5.")
	})

	t.Run("number pointer", func(t *testing.T) {
		msg, err := getMessage(context.Background(), "max", &MessageOptions{
			Attribute: "foo",
			Value:     ptr(1),
			Arguments: []string{"5"},
		})

		assert.NoError(t, err)
		assert.Equal(t, msg, "The foo must not be greater than 5.")
	})

	t.Run("array", func(t *testing.T) {
		msg, err := getMessage(context.Background(), "max", &MessageOptions{
			Attribute: "foo",
			Value:     []int{1},
			Arguments: []string{"5"},
		})

		assert.NoError(t, err)
		assert.Equal(t, msg, "The foo must not have more than 5 items.")
	})
}
