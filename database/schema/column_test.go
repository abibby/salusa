package schema_test

import (
	"testing"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/schema"
	"github.com/abibby/salusa/internal/test"
)

func TestColumnBuilder(t *testing.T) {
	test.QueryTest(t, []test.Case{
		{
			Name:             "column",
			Builder:          schema.NewColumn("foo", dialects.DataTypeInt32),
			ExpectedSQL:      "\"foo\" INTEGER NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Nullable",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Nullable(),
			ExpectedSQL:      "\"foo\" TEXT",
			ExpectedBindings: []any{},
		},
		{
			Name:             "NotNullable",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Nullable().NotNullable(),
			ExpectedSQL:      "\"foo\" TEXT NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Primary",
			Builder:          schema.NewColumn("foo", dialects.DataTypeInt32).Primary(),
			ExpectedSQL:      "\"foo\" INTEGER PRIMARY KEY NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "AutoIncrement",
			Builder:          schema.NewColumn("foo", dialects.DataTypeInt32).AutoIncrement(),
			ExpectedSQL:      "\"foo\" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Default",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Default("bar"),
			ExpectedSQL:      "\"foo\" TEXT NOT NULL DEFAULT 'bar'",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Default Escape",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Default("bar's"),
			ExpectedSQL:      "\"foo\" TEXT NOT NULL DEFAULT 'bar''s'",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Type",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Type(dialects.DataTypeInt32),
			ExpectedSQL:      "\"foo\" INTEGER NOT NULL",
			ExpectedBindings: []any{},
		},
		{
			Name:             "Unique",
			Builder:          schema.NewColumn("foo", dialects.DataTypeString).Unique(),
			ExpectedSQL:      "\"foo\" TEXT NOT NULL UNIQUE",
			ExpectedBindings: []any{},
		},
		{
			Name:             "DefaultCurrentTime",
			Builder:          schema.NewColumn("foo", dialects.DataTypeDateTime).DefaultCurrentTime(),
			ExpectedSQL:      "\"foo\" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP",
			ExpectedBindings: []any{},
		},
	})
}
