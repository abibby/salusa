package generic

import "github.com/abibby/salusa/database/dialects"

// return runQuery(ctx, tx, helpers.Concat(helpers.Raw("DROP TABLE IF EXISTS "), helpers.Identifier(table)))
func (g *Generic) EncodeDropTableQuery(q *dialects.DropTableQuery) (dialects.RawQuery, error) {
	b := newRawQueryBuilder().AddString("DROP TABLE")
	if q.IfExists {
		b.AddString("IF EXISTS")
	}
	return b.AddString(g.core.Identifier(q.Table)).Build()
}
