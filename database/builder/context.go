package builder

import "context"

// WithContext adds a context to the query that will be used when fetching results.
func (b *Builder) WithContext(ctx context.Context) *Builder {
	b = b.Clone()
	b.ctx = ctx
	b.wheres.ctx = ctx
	b.havings.ctx = ctx
	return b
}

// Context returns the context value from the query.
func (b *Builder) Context() context.Context {
	return b.ctx
}

// Context returns the context value from the query.
func (b *ModelBuilder[T]) Context() context.Context {
	return b.subBuilder.Context()
}
