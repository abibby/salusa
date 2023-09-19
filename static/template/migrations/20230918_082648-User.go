package migrations

import (
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/schema"
)

func init() {
	migrations.Add(&migrate.Migration{
		Name: "20230918_082648-User",
		Up: schema.Create("users", func(table *schema.Blueprint) {
			table.Int("id").Primary().AutoIncrement()
			table.String("username")
			table.Blob("password")
		}),
		Down: schema.DropIfExists("users"),
	})
}
