package builder

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

func (b *Builder[T]) ToSQL(d dialects.Dialect) (string, []any, error) {
	return b.subBuilder.ToSQL(d)
}
func (b *SubBuilder) ToSQL(d dialects.Dialect) (string, []any, error) {
	b = b.Clone()
	for _, scope := range b.scopes.allScopes() {
		b = scope.Apply(b)
	}
	return helpers.Result().
		Add(b.selects).
		Add(b.from).
		Add(b.joins).
		Add(b.wheres).
		Add(b.groupBys).
		Add(b.havings).
		Add(b.orderBys).
		Add(b.limit).
		ToSQL(d)
}
