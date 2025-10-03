package pubsub

import (
	"context"
	"io"
)

type PubSub interface {
	Topic(name string) Topic
}

type Topic interface {
	Publish(ctx context.Context, data []byte)
	Subscribe(ctx context.Context) Subscription
}

type Subscription interface {
	io.Closer
	Next() []byte
}
