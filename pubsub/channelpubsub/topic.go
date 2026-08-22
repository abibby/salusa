package channelpubsub

import (
	"context"

	"github.com/abibby/salusa/pubsub"
	"github.com/google/uuid"
)

type Topic struct {
	ch chan pubsub.Message
}

// Close implements [pubsub.Topic].
func (t *Topic) Close() error {
	return nil
}

// Consume implements [pubsub.Topic].
func (t *Topic) Consume(ctx context.Context) (<-chan pubsub.Message, error) {
	return t.ch, nil
}

// Enqueue implements [pubsub.Topic].
func (t *Topic) Enqueue(ctx context.Context, data []byte) error {
	t.ch <- &Message{
		id:   uuid.NewString(),
		data: data,
	}
	return nil
}
