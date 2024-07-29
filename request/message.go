package request

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

type MessageOptions struct {
	Attribute string
	Value     any
	Arguments []string
	Field     reflect.StructField
}

type Message struct {
	Array   string `json:"array"`
	String  string `json:"string"`
	Numeric string `json:"numeric"`
}
type localMessage Message

func (m *Message) UnmarshalJSON(b []byte) error {
	if b[0] == '"' {
		s := ""
		err := json.Unmarshal(b, &s)
		if err != nil {
			return err
		}
		m.Array = s
		m.String = s
		m.Numeric = s
		return nil
	} else {
		return json.Unmarshal(b, (*localMessage)(m))
	}
}

var messages map[string]*Message

//go:embed lang.json
var lang []byte

func init() {
	err := json.Unmarshal(lang, &messages)
	if err != nil {
		panic(fmt.Errorf("could not parse lang.json: %w", err))
	}
}

func getMessage(ctx context.Context, ruleName string, options *MessageOptions) (string, error) {
	if message, ok := options.Field.Tag.Lookup("message"); ok {
		return message, nil
	}
	defaultMessage := func() string {
		if len(options.Arguments) == 0 {
			return ruleName
		}
		return ruleName + " " + strings.Join(options.Arguments, ", ")
	}
	message, ok := messages[ruleName]
	if !ok {
		return defaultMessage(), nil
	}

	messageTemplate := ""

	val := reflect.ValueOf(options.Value)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.CanInt() || val.CanUint() || val.CanFloat() {
		messageTemplate = message.Numeric
	} else if val.Kind() == reflect.String {
		messageTemplate = message.String
	} else if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		messageTemplate = message.Array
	} else {
		messageTemplate = message.String
	}

	t, err := template.New(ruleName).Parse(messageTemplate)
	if err != nil {
		return defaultMessage(), fmt.Errorf("failed to parse template: %w", err)
	}
	buff := &bytes.Buffer{}
	err = t.Execute(buff, options)
	if err != nil {
		return defaultMessage(), fmt.Errorf("failed to execute template: %w", err)
	}
	return buff.String(), nil
}
