package builder

import "context"

// WithContext adds a context to the query that will be used when fetching results.
func (b *ModelBuilder[T]) WithContext(ctx context.Context) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WithContext(ctx)
	return b
}

// From sets the table which the query is targeting.
func (b *ModelBuilder[T]) From(table string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.From(table)
	return b
}

// GroupBy sets the "group by" clause to the query.
func (b *ModelBuilder[T]) GroupBy(columns ...string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.GroupBy(columns...)
	return b
}

// GroupBy adds a "group by" clause to the query.
func (b *ModelBuilder[T]) AddGroupBy(columns ...string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.AddGroupBy(columns...)
	return b
}

// Join adds a join clause to the query.
func (b *ModelBuilder[T]) Join(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Join(table, localColumn, operator, foreignColumn)
	return b
}

// LeftJoin adds a left join clause to the query.
func (b *ModelBuilder[T]) LeftJoin(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.LeftJoin(table, localColumn, operator, foreignColumn)
	return b
}

// RightJoin adds a right join clause to the query.
func (b *ModelBuilder[T]) RightJoin(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.RightJoin(table, localColumn, operator, foreignColumn)
	return b
}

// InnerJoin adds an inner join clause to the query.
func (b *ModelBuilder[T]) InnerJoin(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.InnerJoin(table, localColumn, operator, foreignColumn)
	return b
}

// CrossJoin adds a cross join clause to the query.
func (b *ModelBuilder[T]) CrossJoin(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.CrossJoin(table, localColumn, operator, foreignColumn)
	return b
}

// JoinOn adds a join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) JoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.JoinOn(table, cb)
	return b
}

// LeftJoinOn adds a left join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) LeftJoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.LeftJoinOn(table, cb)
	return b
}

// RightJoinOn adds a right join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) RightJoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.RightJoinOn(table, cb)
	return b
}

// InnerJoinOn adds an inner join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) InnerJoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.InnerJoinOn(table, cb)
	return b
}

// CrossJoinOn adds a cross join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) CrossJoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.CrossJoinOn(table, cb)
	return b
}

// Limit set the maximum number of rows to return.
func (b *ModelBuilder[T]) Limit(limit int) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Limit(limit)
	return b
}

// Offset sets the number of rows to skip before returning the result.
func (b *ModelBuilder[T]) Offset(offset int) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Offset(offset)
	return b
}

// OrderBy adds an order by clause to the query.
func (b *ModelBuilder[T]) OrderBy(column string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrderBy(column)
	return b
}

// OrderByDesc adds a descending order by clause to the query.
func (b *ModelBuilder[T]) OrderByDesc(column string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrderByDesc(column)
	return b
}

// Unordered removes all order by clauses from the query.
func (b *ModelBuilder[T]) Unordered() *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Unordered()
	return b
}

// WithScope adds a local scope to a query.
func (b *ModelBuilder[T]) WithScope(scope *Scope) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WithScope(scope)
	return b
}

// WithoutScope removes the given scope from the local scopes.
func (b *ModelBuilder[T]) WithoutScope(scope *Scope) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WithoutScope(scope)
	return b
}

// WithoutGlobalScope removes a global scope from the query.
func (b *ModelBuilder[T]) WithoutGlobalScope(scope *Scope) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WithoutGlobalScope(scope)
	return b
}

// Select sets the columns to be selected.
func (b *ModelBuilder[T]) Select(columns ...string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Select(columns...)
	return b
}

// AddSelect adds new columns to be selected.
func (b *ModelBuilder[T]) AddSelect(columns ...string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.AddSelect(columns...)
	return b
}

// SelectSubquery sets a subquery to be selected.
func (b *ModelBuilder[T]) SelectSubquery(sb QueryBuilder, as string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.SelectSubquery(sb, as)
	return b
}

// AddSelectSubquery adds a subquery to be selected.
func (b *ModelBuilder[T]) AddSelectSubquery(sb QueryBuilder, as string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.AddSelectSubquery(sb, as)
	return b
}

// SelectFunction sets a column to be selected with a function applied.
func (b *ModelBuilder[T]) SelectFunction(function, column string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.SelectFunction(function, column)
	return b
}

// SelectFunction adds a column to be selected with a function applied.
func (b *ModelBuilder[T]) AddSelectFunction(function, column string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.AddSelectFunction(function, column)
	return b
}

// Distinct forces the query to only return distinct results.
func (b *ModelBuilder[T]) Distinct() *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Distinct()
	return b
}

// Where adds a basic where clause to the query.
func (b *ModelBuilder[T]) Where(column, operator string, value any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Where(column, operator, value)
	return b
}

// Having adds a basic having clause to the query.
func (b *ModelBuilder[T]) Having(column, operator string, value any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Having(column, operator, value)
	return b
}

// OrWhere adds an or where clause to the query
func (b *ModelBuilder[T]) OrWhere(column, operator string, value any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrWhere(column, operator, value)
	return b
}

// OrHaving adds an or having clause to the query
func (b *ModelBuilder[T]) OrHaving(column, operator string, value any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrHaving(column, operator, value)
	return b
}

// WhereColumn adds a where clause to the query comparing two columns.
func (b *ModelBuilder[T]) WhereColumn(column, operator string, valueColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WhereColumn(column, operator, valueColumn)
	return b
}

// HavingColumn adds a having clause to the query comparing two columns.
func (b *ModelBuilder[T]) HavingColumn(column, operator string, valueColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.HavingColumn(column, operator, valueColumn)
	return b
}

// OrWhereColumn adds an or where clause to the query comparing two columns.
func (b *ModelBuilder[T]) OrWhereColumn(column, operator string, valueColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrWhereColumn(column, operator, valueColumn)
	return b
}

// OrHavingColumn adds an or having clause to the query comparing two columns.
func (b *ModelBuilder[T]) OrHavingColumn(column, operator string, valueColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrHavingColumn(column, operator, valueColumn)
	return b
}

// WhereIn adds a where in clause to the query.
func (b *ModelBuilder[T]) WhereIn(column string, values []any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WhereIn(column, values)
	return b
}

// HavingIn adds a having in clause to the query.
func (b *ModelBuilder[T]) HavingIn(column string, values []any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.HavingIn(column, values)
	return b
}

// OrWhereIn adds an or where in clause to the query.
func (b *ModelBuilder[T]) OrWhereIn(column string, values []any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrWhereIn(column, values)
	return b
}

// OrHavingIn adds an or having in clause to the query.
func (b *ModelBuilder[T]) OrHavingIn(column string, values []any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrHavingIn(column, values)
	return b
}

// WhereExists add an exists clause to the query.
func (b *ModelBuilder[T]) WhereExists(query QueryBuilder) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WhereExists(query)
	return b
}

// HavingExists add an exists clause to the query.
func (b *ModelBuilder[T]) HavingExists(query QueryBuilder) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.HavingExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *ModelBuilder[T]) OrWhereExists(query QueryBuilder) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrWhereExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *ModelBuilder[T]) OrHavingExists(query QueryBuilder) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrHavingExists(query)
	return b
}

// WhereSubquery adds a where clause to the query comparing a column and a subquery.
func (b *ModelBuilder[T]) WhereSubquery(subquery QueryBuilder, operator string, value any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WhereSubquery(subquery, operator, value)
	return b
}

// HavingSubquery adds a having clause to the query comparing a column and a subquery.
func (b *ModelBuilder[T]) HavingSubquery(subquery QueryBuilder, operator string, value any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.HavingSubquery(subquery, operator, value)
	return b
}

// OrWhereSubquery adds an or where clause to the query comparing a column and a subquery.
func (b *ModelBuilder[T]) OrWhereSubquery(subquery QueryBuilder, operator string, value any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrWhereSubquery(subquery, operator, value)
	return b
}

// OrHavingSubquery adds an or having clause to the query comparing a column and a subquery.
func (b *ModelBuilder[T]) OrHavingSubquery(subquery QueryBuilder, operator string, value any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrHavingSubquery(subquery, operator, value)
	return b
}

// WhereHas adds a relationship exists condition to the query with where clauses.
func (b *ModelBuilder[T]) WhereHas(relation string, cb func(q *Builder) *Builder) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WhereHas(relation, cb)
	return b
}

// HavingHas adds a relationship exists condition to the query with having clauses.
func (b *ModelBuilder[T]) HavingHas(relation string, cb func(q *Builder) *Builder) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.HavingHas(relation, cb)
	return b
}

// OrWhereHas adds a relationship exists condition to the query with where clauses and an or.
func (b *ModelBuilder[T]) OrWhereHas(relation string, cb func(q *Builder) *Builder) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrWhereHas(relation, cb)
	return b
}

// OrHavingHas adds a relationship exists condition to the query with having clauses and an or.
func (b *ModelBuilder[T]) OrHavingHas(relation string, cb func(q *Builder) *Builder) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrHavingHas(relation, cb)
	return b
}

// WhereRaw adds a raw where clause to the query.
func (b *ModelBuilder[T]) WhereRaw(rawSql string, bindings ...any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.WhereRaw(rawSql, bindings...)
	return b
}

// HavingRaw adds a raw having clause to the query.
func (b *ModelBuilder[T]) HavingRaw(rawSql string, bindings ...any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.HavingRaw(rawSql, bindings...)
	return b
}

// OrWhereRaw adds a raw or where clause to the query.
func (b *ModelBuilder[T]) OrWhereRaw(rawSql string, bindings ...any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrWhereRaw(rawSql, bindings...)
	return b
}

// OrHavingRaw adds a raw or having clause to the query.
func (b *ModelBuilder[T]) OrHavingRaw(rawSql string, bindings ...any) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.OrHavingRaw(rawSql, bindings...)
	return b
}

// And adds a group of conditions to the query
func (b *ModelBuilder[T]) And(cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.And(cb)
	return b
}

// HavingAnd adds a group of conditions to the query
func (b *ModelBuilder[T]) HavingAnd(cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.HavingAnd(cb)
	return b
}

// Or adds a group of conditions to the query with an or
func (b *ModelBuilder[T]) Or(cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Or(cb)
	return b
}

// HavingOr adds a group of conditions to the query with an or
func (b *ModelBuilder[T]) HavingOr(cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.HavingOr(cb)
	return b
}
func (b *ModelBuilder[T]) Dump() *ModelBuilder[T] {
	b = b.Clone()
	b.subBuilder = b.subBuilder.Dump()
	return b
}
