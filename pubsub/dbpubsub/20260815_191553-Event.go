package dbpubsub

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20260815_191553-Event",
		Up: schema.Create("events", func(table *schema.Blueprint) {
			table.DateTime("created_at")
			table.DateTime("updated_at")
			table.Int("id").Primary().AutoIncrement()
			table.String("type")
			table.String("status")
			table.Blob("data")
		}),
		Down: schema.DropIfExists("events"),
	})
}
