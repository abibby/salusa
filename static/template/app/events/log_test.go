package events

import (
	"testing"

	"github.com/abibby/salusa/event"
	"github.com/stretchr/testify/assert"
)

func TestLogEventType(t *testing.T) {
	e := &LogEvent{Message: "hello"}
	assert.Equal(t, event.EventType("template:example-event"), e.Type())
}
