package generic

import "github.com/abibby/salusa/database/dialects"

func (g *Generic) EncodeLiteral(v any) (dialects.SQLResult, error) {
	return dialects.SQLResult{
		Query:    g.core.Binding(),
		Bindings: []any{v},
	}, nil
}

func (g *Generic) EncodeAny(v any) (dialects.SQLResult, error) {
	b := ResultBuilder()
	switch v := v.(type) {
	case dialects.QueryBuilder:
		b.Add(Group(g.EncodeQuery(v.Query())))
	case []dialects.Condition:
		b.Add(Group(g.EncodeConditions(v)))
	case dialects.Column:
		b.Add(g.EncodeColumn(&v))
	case dialects.SQLResult:
		b.Add(v, nil)
	case dialects.RawString:
		b.AddString(string(v))
	default:
		b.Add(g.EncodeLiteral(v))
	}
	return b.Build()
}
