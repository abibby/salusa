package generic

import "github.com/abibby/salusa/database/dialects"

func (g *Generic) EncodeDeleteQuery(q *dialects.DeleteQuery) (dialects.SQLResult, error) {
	return resultBuilder().
		AddString("DELETE FROM").
		AddString(g.core.Identifier(q.Table)).
		Add(g.EncodeWheres(q.Wheres)).
		Build()
}
