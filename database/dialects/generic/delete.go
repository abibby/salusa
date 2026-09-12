package generic

import "gosalusa.com/database/dialects"

func (g *Generic) EncodeDeleteQuery(q *dialects.DeleteQuery) (dialects.RawQuery, error) {
	return newRawQueryBuilder().
		AddString("DELETE FROM").
		AddString(g.core.Identifier(q.Table)).
		Add(g.EncodeWheres(q.Wheres)).
		Build()
}
