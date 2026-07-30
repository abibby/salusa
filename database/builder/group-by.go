package builder

// GroupBy sets the "group by" clause to the query.
func (b *Builder) GroupBy(columns ...string) *Builder {
	b.query.GroupBys = columns
	return b
}

// GroupBy adds a "group by" clause to the query.
func (b *Builder) AddGroupBy(columns ...string) *Builder {
	b.query.GroupBys = append(b.query.GroupBys, columns...)
	return b
}
