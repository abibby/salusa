package generic

import "github.com/abibby/salusa/database/dialects"

func (g *Generic) EncodeUpdateQuery(q *dialects.UpdateQuery) (dialects.SQLResult, error) {
	return resultBuilder().
		AddString("UPDATE").
		AddString(g.core.Identifier(q.Table)).
		Add(g.EncodeUpdateSet(q.Values)).
		Add(g.EncodeWheres(q.Wheres)).
		Build()
}

func (g *Generic) EncodeUpdateSet(values map[string]any) (dialects.SQLResult, error) {
	b := resultBuilder().AddString("SET")
	results := make([]dialects.SQLResult, 0, len(values))
	for k, v := range values {
		r, err := resultBuilder().
			AddString(g.core.Identifier(k)).
			AddString("=").
			Add(g.EncodeAny(v)).
			Build()
		if err != nil {
			return dialects.SQLResult{}, err
		}
		results = append(results, r)
	}
	return b.Add(joinResults(results, ", "), nil).Build()
}
