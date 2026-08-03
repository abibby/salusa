package request

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type messageField struct {
	Foo string `json:"foo" message:"custom message"`
}

func TestMessageUnmarshalJSONString(t *testing.T) {
	m := &Message{}
	err := json.Unmarshal([]byte(`"just a string"`), m)
	assert.NoError(t, err)
	assert.Equal(t, "just a string", m.Array)
	assert.Equal(t, "just a string", m.String)
	assert.Equal(t, "just a string", m.Numeric)
}

func TestMessageUnmarshalJSONError(t *testing.T) {
	m := &Message{}
	err := m.UnmarshalJSON([]byte(`"unterminated`))
	assert.Error(t, err)
}

func TestGetMessageUnknownRule(t *testing.T) {
	msg, err := getMessage(context.Background(), "unknown_rule", &MessageOptions{
		Attribute: "foo",
		Value:     "str",
		Arguments: []string{"a"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "unknown_rule a", msg)
}

func TestGetMessageWithFieldMessage(t *testing.T) {
	sf := reflect.TypeFor[messageField]().Field(0)
	msg, err := getMessage(context.Background(), "max", &MessageOptions{
		Attribute: "foo",
		Value:     "str",
		Arguments: []string{"5"},
		Field:     sf,
	})
	assert.NoError(t, err)
	assert.Equal(t, "custom message", msg)
}

func TestGetMessageBoolValue(t *testing.T) {
	msg, err := getMessage(context.Background(), "max", &MessageOptions{
		Attribute: "foo",
		Value:     true,
		Arguments: []string{"5"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "The foo must not be greater than 5 characters.", msg)
}

func TestGetMessageTemplateParseError(t *testing.T) {
	messages["parse_error_test"] = &Message{
		String:  "{{.Bad",
		Numeric: "{{.Bad",
		Array:   "{{.Bad",
	}
	msg, err := getMessage(context.Background(), "parse_error_test", &MessageOptions{
		Attribute: "foo",
		Value:     "str",
		Arguments: []string{"5"},
	})
	assert.Error(t, err)
	assert.Equal(t, "parse_error_test 5", msg)
}

func TestGetMessageTemplateExecuteError(t *testing.T) {
	msg, err := getMessage(context.Background(), "different", &MessageOptions{
		Attribute: "foo",
		Value:     "str",
		Arguments: []string{"a", "b"},
	})
	assert.Error(t, err)
	assert.Equal(t, "different a, b", msg)
}

func TestValidationErrorError(t *testing.T) {
	assert.Equal(t, "validation error: bad", (ValidationError{"a": []string{"bad"}}).Error())
	assert.Equal(t, "validation error (2)", (ValidationError{"a": []string{"x"}, "b": []string{"y"}}).Error())
}

func TestValidationErrorMergeExisting(t *testing.T) {
	e := ValidationError{"a": []string{"x"}}
	e.Merge(ValidationError{"a": []string{"y"}, "b": []string{"z"}})
	assert.Equal(t, ValidationError{"a": []string{"x", "y"}, "b": []string{"z"}}, e)
}
