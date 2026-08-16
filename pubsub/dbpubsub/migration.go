package dbpubsub

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

var Migration = &migrate.Migration{
	Name: "20260811_065216-Event",
	Up: schema.Create("events", func(table *schema.Blueprint) {
		table.DateTime("created_at")
		table.DateTime("updated_at")
		table.Int("id").Primary().AutoIncrement()
		table.String("status").Default(EventPending)
		table.Blob("data")

		table.Index("idx_job_queue_status_id").AddColumn("status").AddColumn("id")
	}),
	Down: schema.DropIfExists("events"),
}
