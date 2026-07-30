package builder

// From sets the table which the query is targeting.
func (b *Builder) From(table string) *Builder {
	b.query.From = table
	return b
}

// GetTable returns the table the query is targeting
func (b *Builder) GetTable() string {
	return string(b.query.From)
}

// GetTable returns the table the query is targeting
func (b *ModelBuilder[T]) GetTable() string {
	return b.builder.GetTable()
}
