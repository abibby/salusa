package dbpubsub

import (
	"context"
	"time"

	"github.com/abibby/salusa/clog"
	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/pubsub"
	"github.com/davecgh/go-spew/spew"
	"github.com/jmoiron/sqlx"
)

type PubSub struct {
	update database.Update
}

var _ pubsub.PubSub = (*PubSub)(nil)

func New(u database.Update) *PubSub {
	return &PubSub{
		update: u,
	}
}

// Topic implements [pubsub.PubSub].
func (p *PubSub) Topic(name string) pubsub.Topic {
	return &Topic{
		topic:  name,
		update: p.update,
	}
}

type Topic struct {
	topic  string
	update database.Update
}

var _ pubsub.Topic = (*Topic)(nil)

func Register(ctx context.Context) error {
	di.RegisterLazySingletonWith(ctx, func(u database.Update) (pubsub.PubSub, error) {
		return New(u), nil
	})
	return nil
}

// Close implements [pubsub.Queue].
func (t *Topic) Close() error {
	return nil
}

// Consume implements [pubsub.Queue].
func (t *Topic) Consume(ctx context.Context) (<-chan pubsub.Message, error) {
	messages := make(chan pubsub.Message)
	go func() {
		t.listen(ctx, messages)
		close(messages)
	}()
	return messages, nil
}

// Enqueue implements [pubsub.Queue].
func (t *Topic) Enqueue(ctx context.Context, data []byte) error {
	return t.update(func(tx *sqlx.Tx) error {
		return model.Save(tx, &Event{
			Data:   data,
			Topic:  t.topic,
			Status: EventPending,
		})
	})
}

// Pop implements [pubsub.Queue].
func (t *Topic) listen(ctx context.Context, messages chan pubsub.Message) {
	tick := time.Tick(10 * time.Second)

	for {
		err := t.update(func(tx *sqlx.Tx) error {
			events, err := EventQuery().
				Where("status", "=", "pending").
				Where("topic", "=", t.topic).
				OrderBy("id").
				Limit(1).
				ForUpdateSkipLocked().
				UpdateReturning(tx, builder.Updates{
					"status":     "processing",
					"updated_at": time.Now(),
				})
			if err != nil {
				return err
			}
			if len(events) == 0 {
				return nil
			}

			spew.Dump(events)
			messages <- &Message{
				event: events[0],
			}

			return nil
		})
		if err != nil {
			clog.Use(ctx).Warn("failed to process queue", "error", err)
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-tick:
		}
	}
}
