package generic

import (
	"sort"

	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeUpdateQuery(q *dialects.UpdateQuery) (dialects.RawQuery, error) {
	return newRawQueryBuilder().
		AddString("UPDATE").
		AddString(g.core.Identifier(q.Table)).
		Add(g.EncodeUpdateSet(q.Values)).
		Add(g.EncodeWheres(q.Wheres)).
		Build()
}

func (g *Generic) EncodeUpdateSet(values map[string]any) (dialects.RawQuery, error) {
	b := newRawQueryBuilder().AddString("SET")
	results := make([]dialects.RawQuery, 0, len(values))
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		r, err := newRawQueryBuilder().
			AddString(g.core.Identifier(k)).
			AddString("=").
			Add(g.EncodeAny(values[k])).
			Build()
		if err != nil {
			return dialects.RawQuery{}, err
		}
		results = append(results, r)
	}
	return b.Add(joinRawQueries(results, ", "), nil).Build()
}
