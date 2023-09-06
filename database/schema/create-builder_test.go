package schema_test

import (
	"context"
	"testing"

	"github.com/abibby/salusa/database/schema"
	"github.com/abibby/salusa/internal/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	test.QueryTest(t, []test.Case{
		{
			Name:             "create table",
			Builder:          schema.Create("foo", func(table *schema.Blueprint) {}),
			ExpectedSQL:      "CREATE TABLE \"foo\" ();",
			ExpectedBindings: []any{},
		},
		{
			Name: "1 column",
			Builder: schema.Create("foo", func(table *schema.Blueprint) {
				table.String("bar")
			}),
			ExpectedSQL:      "CREATE TABLE \"foo\" (\"bar\" TEXT NOT NULL);",
			ExpectedBindings: []any{},
		},
		{
			Name: "2 columns",
			Builder: schema.Create("foo", func(table *schema.Blueprint) {
				table.Int("id")
				table.String("bar")
			}),
			ExpectedSQL:      "CREATE TABLE \"foo\" (\"id\" INTEGER NOT NULL, \"bar\" TEXT NOT NULL);",
			ExpectedBindings: []any{},
		},
		{
			Name: "primary key",
			Builder: schema.Create("foo", func(table *schema.Blueprint) {
				table.Int("id").Primary()
			}),
			ExpectedSQL:      "CREATE TABLE \"foo\" (\"id\" INTEGER PRIMARY KEY NOT NULL);",
			ExpectedBindings: []any{},
		},
		{
			Name: "composite primary key",
			Builder: schema.Create("foo", func(table *schema.Blueprint) {
				table.Int("id1")
				table.Int("id2")
				table.PrimaryKey("id1", "id2")
			}),
			ExpectedSQL:      "CREATE TABLE \"foo\" (\"id1\" INTEGER NOT NULL, \"id2\" INTEGER NOT NULL, PRIMARY KEY (\"id1\", \"id2\"));",
			ExpectedBindings: []any{},
		},
		{
			Name: "index",
			Builder: schema.Create("foo", func(table *schema.Blueprint) {
				table.Int("id")
				table.String("name")
				table.Index("name_index").AddColumn("name")
			}),
			ExpectedSQL:      "CREATE TABLE \"foo\" (\"id\" INTEGER NOT NULL, \"name\" TEXT NOT NULL); CREATE INDEX IF NOT EXISTS \"name_index\" ON \"foo\" (\"name\");",
			ExpectedBindings: []any{},
		},
		{
			Name: "index",
			Builder: schema.Create("foo", func(table *schema.Blueprint) {
				table.Int("id")
				table.ForeignKey("id", "bar", "foo_id")
			}),
			ExpectedSQL:      "CREATE TABLE \"foo\" (\"id\" INTEGER NOT NULL, CONSTRAINT \"id-bar-foo_id\" FOREIGN KEY (\"id\") REFERENCES \"bar\"(\"foo_id\"));",
			ExpectedBindings: []any{},
		},
	})
}

func TestDefaultSQLite(t *testing.T) {
	test.Run(t, "", func(t *testing.T, tx *sqlx.Tx) {
		c := schema.Create("foo", func(table *schema.Blueprint) {
			table.Int("id").Default(1)
			table.Bool("bool").Default(false)
		})
		err := c.Run(context.Background(), tx)
		assert.NoError(t, err)
	})
}
