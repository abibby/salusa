package generic

import "abibby.com/salusa/database/dialects"

func (g *Generic) EncodeLiteral(v any) (dialects.RawQuery, error) {
	return dialects.RawQuery{
		SQL:      g.core.Binding(),
		Bindings: []any{v},
	}, nil
}

func (g *Generic) EncodeAny(v any) (dialects.RawQuery, error) {
	b := newRawQueryBuilder()
	switch v := v.(type) {
	case dialects.QueryBuilder:
		b.Add(group(g.EncodeSelectQuery(v.Query())))
	case []dialects.Condition:
		b.Add(group(g.EncodeConditions(v)))
	case dialects.Column:
		b.Add(g.EncodeColumn(&v))
	case dialects.RawQuery:
		b.Add(v, nil)
	case dialects.RawString:
		b.AddString(string(v))
	default:
		b.Add(g.EncodeLiteral(v))
	}
	return b.Build()
}
