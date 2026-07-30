package generic

import (
	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeQuery(q *dialects.Query) (dialects.SQLResult, error) {
	b := ResultBuilder().
		Add(g.EncodeSelects(&q.Select)).
		Add(g.EncodeFrom(q.From))

	for _, j := range q.Joins {
		b.Add(g.EncodeJoin(&j))
	}

	return b.
		Add(g.EncodeWheres(q.Wheres)).
		Add(g.EncodeGroupBy(q.GroupBys)).
		Add(g.EncodeHavings(q.Havings)).
		Add(g.EncodeOrderBy(q.OrderBys)).
		Add(g.EncodeLimit(&q.Limit)).
		Build()
}
