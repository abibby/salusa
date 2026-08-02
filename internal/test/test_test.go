package test

import (
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/dialects"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTableNames(t *testing.T) {
	assert.Equal(t, "foos", (&Foo{}).Table())
	assert.Equal(t, "bars", (&Bar{}).Table())
	assert.Equal(t, "foo_soft_deletes", (&FooSoftDelete{}).Table())
}

func TestQueryTest(t *testing.T) {
	QueryTest(t, []Case[dialects.QueryBuilder]{
		{
			Name:             "one select",
			Builder:          builder.From[*Foo]().Select("a"),
			ExpectedSQL:      `SELECT "a" FROM "foos"`,
			ExpectedBindings: []any{},
		},
		{
			Name:             "with where",
			Builder:          builder.From[*Foo]().Select("a").Where("id", "=", 1),
			ExpectedSQL:      `SELECT "a" FROM "foos" WHERE "id" = ?`,
			ExpectedBindings: []any{1},
		},
	})
}

func TestDeleteQueryTest(t *testing.T) {
	DeleteQueryTest(t, []Case[dialects.DeleteQueryBuilder]{
		{
			Name:             "delete",
			Builder:          builder.From[*Foo](),
			ExpectedSQL:      `DELETE FROM "foos"`,
			ExpectedBindings: []any{},
		},
		{
			Name:             "delete with where",
			Builder:          builder.From[*Foo]().Where("id", "=", 3),
			ExpectedSQL:      `DELETE FROM "foos" WHERE "id" = ?`,
			ExpectedBindings: []any{3},
		},
	})
}

func TestUpdateQueryTest(t *testing.T) {
	UpdateQueryTest(t, []Case[*dialects.UpdateQuery]{
		{
			Name: "update",
			Builder: &dialects.UpdateQuery{
				Table:  "foos",
				Values: map[string]any{"name": "x"},
				Wheres: []dialects.Condition{
					{Column: dialects.Column{Column: "id"}, Operator: "=", Value: 1},
				},
			},
			ExpectedSQL:      `UPDATE "foos" SET "name" = ? WHERE "id" = ?`,
			ExpectedBindings: []any{"x", 1},
		},
	})
}

type createTableQB struct{ q *dialects.CreateTableQuery }

func (b *createTableQB) CreateTableQuery() *dialects.CreateTableQuery { return b.q }

type alterTableQB struct{ q *dialects.AlterTableQuery }

func (b *alterTableQB) AlterTableQuery() *dialects.AlterTableQuery { return b.q }

func TestCreateTableTest(t *testing.T) {
	CreateTableTest(t, []Case[dialects.CreateTableQueryBuilder]{
		{
			Name: "create table",
			Builder: &createTableQB{q: &dialects.CreateTableQuery{
				IfNotExists: true,
				Table:       "foos",
				Columns: []dialects.ColumnDefinition{
					{Name: "id", Datatype: dialects.DataTypeInt32, Primary: true},
					{Name: "name", Datatype: dialects.DataTypeString},
				},
			}},
			ExpectedSQL:      `CREATE TABLE IF NOT EXISTS "foos" ("id" INTEGER PRIMARY KEY NOT NULL, "name" TEXT NOT NULL);`,
			ExpectedBindings: []any{},
		},
	})
}

func TestAlterTableTest(t *testing.T) {
	AlterTableTest(t, []Case[dialects.AlterTableQueryBuilder]{
		{
			Name: "alter table",
			Builder: &alterTableQB{q: &dialects.AlterTableQuery{
				Table:         "foos",
				DropColumns:   []string{"old"},
				ModifyColumns: []dialects.ColumnDefinition{{Name: "name", Datatype: dialects.DataTypeString}},
				AddColumns:    []dialects.ColumnDefinition{{Name: "age", Datatype: dialects.DataTypeInt32}},
			}},
			ExpectedSQL:      `ALTER TABLE "foos" DROP COLUMN "old"; ALTER TABLE "foos" MODIFY COLUMN "name" TEXT NOT NULL; ALTER TABLE "foos" ADD "age" INTEGER NOT NULL;`,
			ExpectedBindings: []any{},
		},
	})
}

func TestColumnDefinitionTest(t *testing.T) {
	ColumnDefinitionTest(t, []Case[*dialects.ColumnDefinition]{
		{
			Name: "column definition",
			Builder: &dialects.ColumnDefinition{
				Name:          "id",
				Datatype:      dialects.DataTypeInt32,
				Primary:       true,
				AutoIncrement: true,
				Unique:        true,
			},
			ExpectedSQL:      `"id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE`,
			ExpectedBindings: []any{},
		},
	})
}

func TestRawQueryTest(t *testing.T) {
	RawQueryTest(t, []Case[dialects.QueryBuilder]{
		{
			Name:             "raw",
			Builder:          builder.From[*Foo]().Select("a"),
			ExpectedSQL:      `SELECT "a" FROM "foos"`,
			ExpectedBindings: []any{},
		},
	}, func(d dialects.Dialect, b dialects.QueryBuilder) (dialects.RawQuery, error) {
		return d.EncodeSelectQuery(b.Query())
	})
}

func TestEncoderTest(t *testing.T) {
	d := dialects.New()
	EncoderTest(t, d.EncodeSelectQuery, []EncoderTestCase[*dialects.SelectQuery]{
		{
			Name: "select",
			Builder: &dialects.SelectQuery{
				Select: dialects.Select{Columns: []dialects.Column{{Column: "a"}}},
				From:   "foos",
			},
			ExpectedSQL:      `SELECT "a" FROM "foos"`,
			ExpectedBindings: []any{},
		},
	})
}

func TestRunnerRunNoTx(t *testing.T) {
	RunNoTx(t, "select count", func(t *testing.T, db *sqlx.DB) {
		var n int
		err := db.Get(&n, "SELECT COUNT(*) FROM foos")
		assert.NoError(t, err)
	})
}

func TestRunnerRun(t *testing.T) {
	Run(t, "select count", func(t *testing.T, tx *sqlx.Tx) {
		var n int
		err := tx.Get(&n, "SELECT COUNT(*) FROM foos")
		assert.NoError(t, err)
	})
}
