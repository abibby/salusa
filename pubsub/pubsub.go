package pubsub

import (
	"context"
	"io"

	"github.com/abibby/salusa/di"
)

type PubSub interface {
	Topic(name string) Topic
}

type Topic interface {
	Producer
	Consumer
}

type Subscription interface {
	io.Closer
	Next() ([]byte, error)
}

type Message interface {
	ID() string
	Data() []byte

	// Ack signals the message was processed successfully and can be deleted.
	Ack(ctx context.Context) error

	// Nack signals processing failed. The message should be re-queued or dead-lettered.
	Nack(ctx context.Context) error
}

type Producer interface {
	Enqueue(ctx context.Context, data []byte) error
	Close() error
}

type Consumer interface {
	// Consume returns a read-only channel of messages.
	// The backend handles batching, pre-fetching, and pushing to this channel.
	Consume(ctx context.Context) (<-chan Message, error)
	Close() error
}

func RegisterTopic(ctx context.Context) {
	di.RegisterWith(ctx, func(ctx context.Context, tag string, with PubSub) (Topic, error) {
		return with.Topic(tag), nil
	})
}
