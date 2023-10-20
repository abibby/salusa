package event

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"log/slog"
	"reflect"

	"github.com/abibby/salusa/clog"
)

type EventType string

type Event interface {
	Type() EventType
	WithContext(ctx context.Context)
	Context(ctx context.Context) context.Context
}

type EventLogger struct {
	Logger *slog.Logger
}

func (e *EventLogger) Context(ctx context.Context) context.Context {
	if e.Logger == nil {
		return ctx
	}
	return clog.Update(ctx, func(l *slog.Logger) *slog.Logger {
		return e.Logger
	})
}
func (e *EventLogger) WithContext(ctx context.Context) {
	e.Logger = clog.Use(ctx)
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
