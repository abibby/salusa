package generic

import "github.com/abibby/salusa/database/dialects"

func (g *Generic) EncodeQuery(q *dialects.Query) (dialects.SQLResult, error) {
	// return helpers.Result().
	// 	Add(b.selects).
	// 	Add(b.from).
	// 	Add(b.joins).
	// 	Add(b.wheres).
	// 	Add(b.groupBys).
	// 	Add(b.havings).
	// 	Add(b.orderBys).
	// 	Add(b.limit).
	// 	SQLString(d)
	// JoinResults()

	return ResultBuilder().
		Add(g.EncodeSelects(&q.Select)).
		Add(g.EncodeFrom(q.From)).
		Add(g.EncodeWheres(q.Wheres)).
		Add(g.EncodeHavings(q.Havings)).
		Build()
}
