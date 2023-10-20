package event

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestEvent1 struct {
	EventLogger
	Foo string
}

func (e *TestEvent1) Type() EventType {
	return "test-event:1"
}

type TestEvent2 struct {
	EventLogger
	Bar string
}

func (e *TestEvent2) Type() EventType {
	return "test-event:2"
}

func TestEncodeDecode(t *testing.T) {
	t2 := &TestEvent2{
		Bar: "baz",
	}
	b, err := encodeEvent(t2)
	if !assert.NoError(t, err) {
		return
	}

	events := map[EventType]reflect.Type{
		(&TestEvent1{}).Type(): reflect.TypeOf(&TestEvent1{}),
		(&TestEvent2{}).Type(): reflect.TypeOf(&TestEvent2{}),
	}
	e, err := decodeEvent(b, events)
	assert.NoError(t, err)
	assert.IsType(t, &TestEvent2{}, e)
	assert.Equal(t, "baz", e.(*TestEvent2).Bar)
}
