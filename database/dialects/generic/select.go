package generic

import (
	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeSelects(s *dialects.Select) (dialects.SQLResult, error) {
	if len(s.Columns) == 0 {
		return dialects.SQLResult{
			Bindings: []any{},
		}, nil
	}

	b := ResultBuilder()
	b.AddString("SELECT")
	if s.Distinct {
		b.AddString("DISTINCT")
	}

	columns := make([]dialects.SQLResult, len(s.Columns))
	var err error
	for i, c := range s.Columns {
		columns[i], err = g.EncodeColumn(&c)
		if err != nil {
			return dialects.SQLResult{}, err
		}
	}
	b.Add(JoinResults(columns, ", "), nil)

	return b.Build()
}

func (g *Generic) EncodeFunctionCall(fc *dialects.FunctionCall) (dialects.SQLResult, error) {
	return ResultBuilder().
		AddString(fc.Name + "(" + g.core.Identifier(fc.Arguments) + ")").
		Build()
}

func (g *Generic) EncodeColumn(c *dialects.Column) (dialects.SQLResult, error) {
	b := ResultBuilder()
	if c.Function != nil {
		b.Add(g.EncodeFunctionCall(c.Function))
	} else {
		b.AddString(g.core.Identifier(c.Column))
	}
	if c.As != "" {
		b.AddString("AS").AddString(g.core.Identifier(c.As))
	}
	return b.Build()
}

// func columnIdentifier(column string) helpers.SQLStringer {
// 	var identifier helpers.SQLStringer
// 	parts := strings.SplitN(column, " as ", 2)
// 	c := parts[0]
// 	if c == "*" {
// 		identifier = helpers.Raw("*")
// 	} else if strings.HasSuffix(c, ".*") {
// 		identifier = helpers.Concat(helpers.Identifier(c[:len(c)-2]), helpers.Raw(".*"))
// 	} else {
// 		identifier = helpers.Identifier(c)
// 	}
// 	if len(parts) > 1 {
// 		identifier = helpers.Concat(
// 			identifier,
// 			helpers.Raw(" as "),
// 			helpers.Identifier(parts[1]),
// 		)
// 	}
// 	return identifier
// }
