package dbpubsub

import (
	"context"
	"strconv"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/pubsub"
	"github.com/jmoiron/sqlx"
)

type Message struct {
	event  *Event
	update database.Update
}

var _ pubsub.Message = (*Message)(nil)

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
	return m.update(func(tx *sqlx.Tx) error {
		m.event.Status = EventFinished
		return model.SaveContext(ctx, tx, m.event)
	})
}

// Nack implements [event.Message].
func (m *Message) Nack(ctx context.Context) error {
	return m.update(func(tx *sqlx.Tx) error {
		m.event.Status = EventError
		return model.SaveContext(ctx, tx, m.event)
	})
}
