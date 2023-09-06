package requests

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"text/template"
)

type MessageOptions struct {
	Attribute string
	Value     any
	Arguments []string
}

type Message struct {
	Array   string `json:"array"`
	String  string `json:"string"`
	Numeric string `json:"numeric"`
}

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
		type localMessage Message
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

func getMessage(ruleName string, options *MessageOptions) string {
	defaultMessage := func() string {
		if len(options.Arguments) == 0 {
			return ruleName
		}
		return ruleName + " " + strings.Join(options.Arguments, ", ")
	}
	message, ok := messages[ruleName]
	if !ok {
		return defaultMessage()
	}

	messageTemplate := ""

	switch options.Value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		messageTemplate = message.Numeric
	case string:
		messageTemplate = message.String
	default:
		messageTemplate = message.String
	}

	t, err := template.New(ruleName).Parse(messageTemplate)
	if err != nil {
		log.Print(fmt.Errorf("failed to parse template: %w", err))
		return defaultMessage()
	}
	buff := &bytes.Buffer{}
	err = t.Execute(buff, options)
	if err != nil {
		log.Print(fmt.Errorf("failed to execute template: %w", err))
		return defaultMessage()
	}
	return buff.String()
}
