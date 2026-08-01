package generic

import (
	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeSelects(s *dialects.Select) (dialects.RawQuery, error) {
	if len(s.Columns) == 0 {
		return dialects.RawQuery{
			Bindings: []any{},
		}, nil
	}

	b := resultBuilder()
	b.AddString("SELECT")
	if s.Distinct {
		b.AddString("DISTINCT")
	}

	columns := make([]dialects.RawQuery, len(s.Columns))
	var err error
	for i, c := range s.Columns {
		columns[i], err = g.EncodeColumn(&c)
		if err != nil {
			return dialects.RawQuery{}, err
		}
	}
	b.Add(joinResults(columns, ", "), nil)

	return b.Build()
}

func (g *Generic) EncodeFunctionCall(fc *dialects.FunctionCall) (dialects.RawQuery, error) {
	return resultBuilder().
		AddString(fc.Name + "(" + g.core.Identifier(fc.Arguments) + ")").
		Build()
}

func (g *Generic) EncodeColumn(c *dialects.Column) (dialects.RawQuery, error) {
	b := resultBuilder()
	if c.Column != "" {
		b.AddString(g.core.Identifier(c.Column))
	} else if c.Function != nil {
		b.Add(g.EncodeFunctionCall(c.Function))
	} else if c.SubQuery != nil {
		b.Add(group(g.EncodeSelectQuery(c.SubQuery.Query())))
	} else {
	}
	if c.As != "" {
		b.AddString("AS").AddString(g.core.Identifier(c.As))
	}
	return b.Build()
}
