package generic

import "github.com/abibby/salusa/database/dialects"

func (g *Generic) EncodeFrom(from string) (dialects.RawQuery, error) {
	if from == "" {
		return dialects.RawQuery{}, nil
	}
	return dialects.RawQuery{
		Query: "FROM " + g.core.Identifier(from),
	}, nil
}
