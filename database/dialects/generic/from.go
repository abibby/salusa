package generic

import "abibby.com/salusa/database/dialects"

func (g *Generic) EncodeFrom(from string) (dialects.RawQuery, error) {
	if from == "" {
		return dialects.RawQuery{}, nil
	}
	return dialects.RawQuery{
		SQL: "FROM " + g.core.Identifier(from),
	}, nil
}
