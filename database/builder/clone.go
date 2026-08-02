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
		query:   *b.query.Clone(),
		scopes:  b.scopes.Clone(),
		wheres:  b.wheres.Clone(),
		havings: b.havings.Clone(),
		ctx:     b.ctx,
	}
}

func cloneSlice[T any](arr []T) []T {
	l := make([]T, len(arr))
	copy(l, arr)
	return l
}
