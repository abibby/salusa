package builder

func (b *Builder[T]) With(withs ...string) *Builder[T] {
	b.withs = append(b.withs, withs...)
	return b
}
