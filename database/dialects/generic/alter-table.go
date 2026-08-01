package generic

import "github.com/abibby/salusa/database/dialects"

func (g *Generic) EncodeAlterTableQuery(q *dialects.AlterTableQuery) (dialects.SQLResult, error) {
	b := resultBuilder()
	return b.Build()
}
