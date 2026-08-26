package generic_test

import (
	"testing"

	"abibby.com/salusa/database/dialects"
	"abibby.com/salusa/database/dialects/generic"
	"abibby.com/salusa/internal/test"
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
			ExpectedSQL:      "`table`.`column`",
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
			ExpectedSQL:      "DISTINCT `table`.`column`",
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
			ExpectedSQL:      "`table1`.`column1`, `table2`.`column2`",
			ExpectedBindings: []any{},
		},
		{
			Name: "as",
			Builder: &dialects.Select{
				Columns: []dialects.Column{
					{Column: "table.column", As: "name"},
				},
			},
			ExpectedSQL:      "`table`.`column` AS `name`",
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
			ExpectedSQL:      "count(*)",
			ExpectedBindings: []any{},
		},
		{
			Name: "disinct",
			Builder: &dialects.Select{
				Columns: []dialects.Column{
					{Column: "foo"},
				},
				Distinct: true,
			},
			ExpectedSQL:      "DISTINCT `foo`",
			ExpectedBindings: []any{},
		},
	})
}
