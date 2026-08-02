package generic

import (
	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeJoins(joins []dialects.Join) (dialects.RawQuery, error) {
	b := newRawQueryBuilder()
	for _, j := range joins {
		b.Add(g.EncodeJoin(&j))
	}
	return b.Build()
}
func (g *Generic) EncodeJoin(j *dialects.Join) (dialects.RawQuery, error) {
	b := newRawQueryBuilder().
		AddString(j.Direction).
		AddString("JOIN").
		AddString(g.core.Identifier(j.Table))

	if len(j.Conditions) > 0 {
		b.AddString("ON").Add(g.EncodeConditions(j.Conditions))
	}

	return b.Build()
}
