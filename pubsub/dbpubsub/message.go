package dbpubsub

import (
	"context"
	"strconv"

	"github.com/abibby/salusa/event"
)

type Message struct {
	event *Event
}

var _ event.Message = (*Message)(nil)

// Data implements [event.Message].
func (m *Message) Data() []byte {
	return m.event.Data
}

// ID implements [event.Message].
func (m *Message) ID() string {
	return strconv.FormatInt(int64(m.event.ID), 10)
}

// Ack implements [event.Message].
func (m *Message) Ack(ctx context.Context) error {
	panic("unimplemented")
}

// Nack implements [event.Message].
func (m *Message) Nack(ctx context.Context) error {
	panic("unimplemented")
}
