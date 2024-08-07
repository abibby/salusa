package builder

func (b *ModelBuilder[T]) Clone() *ModelBuilder[T] {
	return &ModelBuilder[T]{
		builder:       b.builder.Clone(),
		withs:         cloneSlice(b.withs),
		withoutScopes: b.withoutScopes.Clone(),
	}
}
func (b *Builder) Clone() *Builder {
	return &Builder{
		selects:  b.selects.Clone(),
		from:     b.from.Clone(),
		joins:    b.joins.Clone(),
		wheres:   b.wheres.Clone(),
		groupBys: b.groupBys.Clone(),
		havings:  b.havings.Clone(),
		limit:    b.limit.Clone(),
		orderBys: b.orderBys.Clone(),
		scopes:   b.scopes.Clone(),
		ctx:      b.ctx,
	}
}

func cloneSlice[T any](arr []T) []T {
	l := make([]T, len(arr))
	copy(l, arr)
	return l
}
