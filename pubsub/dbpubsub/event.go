package dbpubsub

import (
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/database/model/mixins"
)

//go:generate spice generate:migration
type Event struct {
	model.BaseModel
	mixins.Timestamps

	ID     int         `json:"id"     db:"id,primary,autoincrement"`
	Status EventStatus `json:"status" db:"status"`
	Topic  string      `json:"topic"  db:"topic"`
	Data   []byte      `json:"data"   db:"data"`
}

type EventStatus string

var (
	EventPending    = EventStatus("pending")
	EventProcessing = EventStatus("processing")
	EventFinished   = EventStatus("finished")
	EventError      = EventStatus("error")
)

func EventQuery() *builder.ModelBuilder[*Event] {
	return builder.From[*Event]()
}
