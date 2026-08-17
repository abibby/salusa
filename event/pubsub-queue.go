package event

import (
	"context"
	"reflect"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/pubsub"
)

type PubSubQueue struct {
	topic    pubsub.Topic
	messages <-chan pubsub.Message
}

func NewPubSubQueue(topic pubsub.Topic) *PubSubQueue {
	return &PubSubQueue{
		topic: topic,
	}
}

func (q *PubSubQueue) Push(ctx context.Context, e Event) error {
	b, err := encodeEvent(e)
	if err != nil {
		return err
	}
	q.topic.Enqueue(ctx, b)
	return nil
}
func (q *PubSubQueue) Pop(ctx context.Context, events map[EventType]reflect.Type) (Event, error) {
	if q.messages == nil {
		messages, err := q.topic.Consume(ctx)
		if err != nil {
			return nil, err
		}
		q.messages = messages
	}
	msg := <-q.messages
	msg.Ack(ctx)
	return decodeEvent(msg.Data(), events)
}

func RegisterPubSubQueue(ctx context.Context) error {
	di.RegisterWith(ctx, func(ctx context.Context, tag string, ps pubsub.PubSub) (*PubSubQueue, error) {
		return NewPubSubQueue(ps.Topic(tag)), nil
	})

	return nil
}
