package event

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
)

type EventType string

type Event interface {
	Type() EventType
}

var (
	ErrEventTypeNotFound = errors.New("event type not found")
)

func encodeEvent(e Event) ([]byte, error) {
	buff := bytes.NewBufferString(string(e.Type()) + "|")
	err := gob.NewEncoder(buff).Encode(e)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func decodeEvent(b []byte, events map[EventType]reflect.Type) (Event, error) {
	parts := bytes.SplitN(b, []byte{'|'}, 2)
	eventType := EventType(parts[0])
	data := parts[1]

	t, ok := events[eventType]
	if !ok {
		return nil, ErrEventTypeNotFound
	}

	dereferenced := false
	if t.Kind() == reflect.Ptr {
		dereferenced = true
		t = t.Elem()
	}
	v := reflect.New(t)
	buff := bytes.NewBuffer(data)
	err := gob.NewDecoder(buff).DecodeValue(v)
	if err != nil {
		return nil, err
	}
	if !dereferenced {
		v = v.Elem()
	}
	return v.Interface().(Event), nil
}
