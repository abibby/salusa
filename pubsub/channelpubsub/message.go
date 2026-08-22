package channelpubsub

import (
	"context"

	"github.com/abibby/salusa/pubsub"
)

type Message struct {
	id   string
	data []byte
}

var _ pubsub.Message = (*Message)(nil)

// Data implements [pubsub.Message].
func (m *Message) Data() []byte {
	return m.data
}

// ID implements [pubsub.Message].
func (m *Message) ID() string {
	return m.id
}

// Ack implements [pubsub.Message].
func (m *Message) Ack(ctx context.Context) error {
	return nil
}

// Nack implements [pubsub.Message].
func (m *Message) Nack(ctx context.Context) error {
	return nil
}
