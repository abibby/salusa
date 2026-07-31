package generic

import (
	"strings"

	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeOrderBy(orderBys []string) (dialects.SQLResult, error) {
	if len(orderBys) == 0 {
		return dialects.SQLResult{}, nil
	}
	b := resultBuilder()
	b.AddString("ORDER BY")

	identifiers := make([]string, len(orderBys))
	for i, group := range orderBys {
		identifiers[i] = g.core.Identifier(group)
	}
	return b.AddString(strings.Join(identifiers, ", ")).Build()
}
