package dbpubsub

import (
	"time"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/database/model/mixins"
)

//go:generate spice generate:migration
type Event struct {
	model.BaseModel
	mixins.Timestamps

	ID       int         `db:"id,primary,autoincrement"`
	FirstRun *time.Time  `db:"first_run"`
	RunAt    *time.Time  `db:"run_at"`
	Status   EventStatus `db:"status"`
	Topic    string      `db:"topic"`
	Data     []byte      `db:"data"`
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
