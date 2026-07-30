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
