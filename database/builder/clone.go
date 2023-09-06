package builder

func (b *Builder[T]) Clone() *Builder[T] {
	return &Builder[T]{
		subBuilder:    b.subBuilder.Clone(),
		withs:         cloneSlice(b.withs),
		withoutScopes: b.withoutScopes.Clone(),
	}
}
func (b *SubBuilder) Clone() *SubBuilder {
	return &SubBuilder{
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
	for i, v := range arr {
		l[i] = v
	}
	return l
}
