package event

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type valEvent struct {
	Foo string
}

func (valEvent) Type() EventType {
	return "val-event"
}

func TestDecodeEventErrors(t *testing.T) {
	t.Run("event type not found", func(t *testing.T) {
		_, err := decodeEvent([]byte("unknown|data"), map[EventType]reflect.Type{})
		assert.ErrorIs(t, err, ErrEventTypeNotFound)
	})

	t.Run("bad gob data", func(t *testing.T) {
		_, err := decodeEvent([]byte("test-event:1|notgobdata"), map[EventType]reflect.Type{
			(&TestEvent1{}).Type(): reflect.TypeOf(&TestEvent1{}),
		})
		assert.Error(t, err)
	})

	t.Run("value type event", func(t *testing.T) {
		b, err := encodeEvent(valEvent{Foo: "bar"})
		assert.NoError(t, err)

		e, err := decodeEvent(b, map[EventType]reflect.Type{
			"val-event": reflect.TypeOf(valEvent{}),
		})
		assert.NoError(t, err)
		val, ok := e.(valEvent)
		assert.True(t, ok)
		assert.Equal(t, "bar", val.Foo)
	})
}

type badFuncEvent struct {
	Fn func()
}

func (*badFuncEvent) Type() EventType {
	return "bad-event"
}

func TestEncodeEventError(t *testing.T) {
	_, err := encodeEvent(&badFuncEvent{})
	assert.Error(t, err)
}

type TestEvent1 struct {
	Foo string
}

func (e *TestEvent1) Type() EventType {
	return "test-event:1"
}

type TestEvent2 struct {
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
