package builder

// Limit set the maximum number of rows to return.
func (b *Builder) Limit(limit int) *Builder {
	b.query.Limit.Limit = limit
	return b
}

// Offset sets the number of rows to skip before returning the result.
func (b *Builder) Offset(offset int) *Builder {
	b.query.Limit.Offset = offset
	return b
}
