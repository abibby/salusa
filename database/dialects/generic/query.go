package generic

import (
	"github.com/abibby/salusa/database/dialects"
)

func (g *Generic) EncodeSelectQuery(q *dialects.SelectQuery) (dialects.RawQuery, error) {
	return newRawQueryBuilder().
		AddString("SELECT").
		Add(g.EncodeSelects(&q.Select)).
		Add(g.EncodeFrom(q.From)).
		Add(g.EncodeJoins(q.Joins)).
		Add(g.EncodeWheres(q.Wheres)).
		Add(g.EncodeGroupBy(q.GroupBys)).
		Add(g.EncodeHavings(q.Havings)).
		Add(g.EncodeOrderBy(q.OrderBys)).
		Add(g.EncodeLimit(&q.Limit)).
		Build()
}
