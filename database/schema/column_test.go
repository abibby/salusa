package schema_test

import (
	"testing"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/schema"
	"github.com/abibby/salusa/internal/test"
)

func TestColumnBuilder(t *testing.T) {
	test.ColumnDefinitionTest(t, []test.Case[*dialects.ColumnDefinition]{
		{
			Name:             "column",
			Builder:          schema.NewColumn("foo", dialects.DataTypeInt32).ColumnDefinition(),
			ExpectedSQL:      "\"foo\" INTEGER NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Nullable",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Nullable().ColumnDefinition(),
			ExpectedSQL:      "\"foo\" TEXT",
			ExpectedBindings: []any{},
		},
		{
			Name:             "NotNullable",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Nullable().NotNullable().ColumnDefinition(),
			ExpectedSQL:      "\"foo\" TEXT NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Primary",
			Builder:          schema.NewColumn("foo", dialects.DataTypeInt32).Primary().ColumnDefinition(),
			ExpectedSQL:      "\"foo\" INTEGER PRIMARY KEY NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "AutoIncrement",
			Builder:          schema.NewColumn("foo", dialects.DataTypeInt32).AutoIncrement().ColumnDefinition(),
			ExpectedSQL:      "\"foo\" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Default",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Default("bar").ColumnDefinition(),
			ExpectedSQL:      "\"foo\" TEXT NOT NULL DEFAULT 'bar'",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Default Escape",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Default("bar's").ColumnDefinition(),
			ExpectedSQL:      "\"foo\" TEXT NOT NULL DEFAULT 'bar''s'",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Type",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Type(dialects.DataTypeInt32).ColumnDefinition(),
			ExpectedSQL:      "\"foo\" INTEGER NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Unique",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Unique().ColumnDefinition(),
			ExpectedSQL:      "\"foo\" TEXT NOT NULL UNIQUE",
			ExpectedBindings: []any{},
		},
		{
			Name:             "DefaultCurrentTime",
			Builder:          schema.NewColumn("foo", dialects.DataTypeDateTime).DefaultCurrentTime().ColumnDefinition(),
			ExpectedSQL:      "\"foo\" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP",
			ExpectedBindings: []any{},
		},
	})
}
