package generic_test

import (
	"testing"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/generic"
	"github.com/abibby/salusa/internal/test"
)

func TestGeneric_EncodeCreateTableQuery(t *testing.T) {
	g := generic.New(&testCore{})
	test.EncoderTest(t, g.EncodeCreateTableQuery, []test.EncoderTestCase[*dialects.CreateTableQuery]{
		{
			Name: "simple",
			Builder: &dialects.CreateTableQuery{
				Table: "foo",
				Columns: []dialects.ColumnDefinition{
					{Name: "c1", Datatype: dialects.DataTypeDate},
				},
			},
			ExpectedSQL:      "CREATE TABLE `foo` (\n\t`c1` date NOT NULL\n)",
			ExpectedBindings: []any{},
		},
		{
			Name: "multi",
			Builder: &dialects.CreateTableQuery{
				Table: "foo",
				Columns: []dialects.ColumnDefinition{
					{Name: "c1", Datatype: dialects.DataTypeDate},
					{Name: "c2", Datatype: dialects.DataTypeInt64},
				},
			},
			ExpectedSQL: "CREATE TABLE `foo` (\n" +
				"\t`c1` date NOT NULL,\n" +
				"\t`c2` int64 NOT NULL\n" +
				")",
			ExpectedBindings: []any{},
		},
		{
			Name: "full",
			Builder: &dialects.CreateTableQuery{
				Table: "foo",
				Columns: []dialects.ColumnDefinition{
					{Name: "id", Datatype: dialects.DataTypeInt32, Primary: true},
					{Name: "name", Datatype: dialects.DataTypeString, Unique: true},
					{Name: "optional", Datatype: dialects.DataTypeString, Nullable: true},
				},
			},
			ExpectedSQL: "CREATE TABLE `foo` (\n" +
				"\t`id` int32 PRIMARY KEY NOT NULL,\n" +
				"\t`name` string NOT NULL UNIQUE,\n" +
				"\t`optional` string\n" +
				")",
			ExpectedBindings: []any{},
		},
	})
}
