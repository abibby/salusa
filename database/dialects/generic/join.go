package generic

import (
	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeJoin(j *dialects.Join) (dialects.SQLResult, error) {
	b := ResultBuilder().
		AddString(j.Direction).
		AddString("JOIN").
		AddString(g.core.Identifier(j.Table))

	if len(j.Conditions) > 0 {
		b.AddString("ON").Add(g.EncodeConditions(j.Conditions))
	}

	return b.Build()
}
