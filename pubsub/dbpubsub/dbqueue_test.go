package dbpubsub_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/internal/test"
	"github.com/abibby/salusa/pubsub/dbpubsub"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type valEvent struct {
	Foo string
}

func (valEvent) Type() event.EventType {
	return "val-event"
}

func TestQueue_Push_Pop(t *testing.T) {
	t.SkipNow()
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		err := dbpubsub.Migration.Up.Run(t.Context(), tx)
		if !assert.NoError(t, err) {
			return
		}
		p := dbpubsub.New(database.Update(func(f func(tx *sqlx.Tx) error) error {
			return f(tx)
		}))
		ctx, cancel := context.WithCancel(context.Background())

		q := p.Topic("default")

		q.Enqueue(ctx, []byte("a"))

		messages, err := q.Consume(ctx)
		if assert.NoError(t, err) {
			m := <-messages
			assert.Equal(t, []byte("a"), m.Data())

		}
		cancel()
	})
}
