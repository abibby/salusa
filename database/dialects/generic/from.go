package generic

import "github.com/abibby/salusa/database/dialects"

func (g *Generic) EncodeFrom(from string) (dialects.SQLResult, error) {
	if from == "" {
		return dialects.SQLResult{}, nil
	}
	return dialects.SQLResult{
		Query: "FROM " + g.core.Identifier(from),
	}, nil
}

func (g *Generic) EncodeLiteral(v any) (dialects.SQLResult, error) {
	return dialects.SQLResult{
		Query:    g.core.Binding(),
		Bindings: []any{v},
	}, nil
}
