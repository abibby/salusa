package builder

// OrderBy adds an order by clause to the query.
func (b *Builder) OrderBy(column string) *Builder {
	// return append(o, helpers.Identifier(column))
	b.query.OrderBys = append(b.query.OrderBys, column)
	return b
}

// OrderByDesc adds a descending order by clause to the query.
func (b *Builder) OrderByDesc(column string) *Builder {
	// return append(o, helpers.Join([]helpers.SQLStringer{helpers.Identifier(column), helpers.Raw("DESC")}, " "))
	// b.query.OrderBys = append(b.query.OrderBys, column)
	return b
}

// Unordered removes all order by clauses from the query.
func (b *Builder) Unordered() *Builder {
	// return orderBys{}
	b.query.OrderBys = []string{}
	return b
}
