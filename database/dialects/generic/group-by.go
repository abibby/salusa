package generic

import (
	"strings"

	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeGroupBy(groups []string) (dialects.RawQuery, error) {
	if len(groups) == 0 {
		return dialects.RawQuery{}, nil
	}
	b := newRawQueryBuilder()
	b.AddString("GROUP BY")

	identifiers := make([]string, len(groups))
	for i, group := range groups {
		identifiers[i] = g.core.Identifier(group)
	}
	return b.AddString(strings.Join(identifiers, ", ")).Build()
}
