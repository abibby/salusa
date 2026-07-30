package generic_test

import (
	"testing"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/generic"
	"github.com/abibby/salusa/internal/test"
)

func TestGeneric_EncodeSelects(t *testing.T) {
	g := generic.New(&testCore{})
	test.EncoderTest(t, g.EncodeSelects, []test.EncoderTestCase[*dialects.Select]{
		{
			Name:             "empty",
			Builder:          &dialects.Select{},
			ExpectedSQL:      "",
			ExpectedBindings: []any{},
		},
		{
			Name: "single",
			Builder: &dialects.Select{
				Columns: []dialects.Column{
					{Column: "table.column"},
				},
			},
			ExpectedSQL:      "SELECT `table`.`column`",
			ExpectedBindings: []any{},
		},
		{
			Name: "distinct",
			Builder: &dialects.Select{
				Distinct: true,
				Columns: []dialects.Column{
					{Column: "table.column"},
				},
			},
			ExpectedSQL:      "SELECT DISTINCT `table`.`column`",
			ExpectedBindings: []any{},
		},
		{
			Name: "multi",
			Builder: &dialects.Select{
				Columns: []dialects.Column{
					{Column: "table1.column1"},
					{Column: "table2.column2"},
				},
			},
			ExpectedSQL:      "SELECT `table1`.`column1`, `table2`.`column2`",
			ExpectedBindings: []any{},
		},
		{
			Name: "as",
			Builder: &dialects.Select{
				Columns: []dialects.Column{
					{Column: "table.column", As: "name"},
				},
			},
			ExpectedSQL:      "SELECT `table`.`column` AS `name`",
			ExpectedBindings: []any{},
		},
		{
			Name: "function",
			Builder: &dialects.Select{
				Columns: []dialects.Column{
					{Function: &dialects.FunctionCall{
						Name:      "count",
						Arguments: "*",
					}},
				},
			},
			ExpectedSQL:      "SELECT count(*)",
			ExpectedBindings: []any{},
		},
	})
}
