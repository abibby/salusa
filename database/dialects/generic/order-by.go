package generic

import (
	"strings"

	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeOrderBy(orderBys []dialects.OrderColumn) (dialects.SQLResult, error) {
	if len(orderBys) == 0 {
		return dialects.SQLResult{}, nil
	}

	identifiers := make([]string, len(orderBys))
	for i, group := range orderBys {
		identifiers[i] = g.core.Identifier(group.Column)
		if group.Descending {
			identifiers[i] += " DESC"
		}
	}
	return dialects.Raw("ORDER BY " + strings.Join(identifiers, ", ")), nil
}
