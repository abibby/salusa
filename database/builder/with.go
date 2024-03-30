package builder

func (b *ModelBuilder[T]) With(withs ...string) *ModelBuilder[T] {
	b.withs = append(b.withs, withs...)
	return b
}
