package builder

// From sets the table which the query is targeting.
func (b *Builder) From(table string) *Builder {
	b = b.Clone()
	b.from = b.from.From(table)
	return b
}

// GroupBy sets the "group by" clause to the query.
func (b *Builder) GroupBy(columns ...string) *Builder {
	b = b.Clone()
	b.groupBys = b.groupBys.GroupBy(columns...)
	return b
}

// GroupBy adds a "group by" clause to the query.
func (b *Builder) AddGroupBy(columns ...string) *Builder {
	b = b.Clone()
	b.groupBys = b.groupBys.AddGroupBy(columns...)
	return b
}

// Join adds a join clause to the query.
func (b *Builder) Join(table, localColumn, operator, foreignColumn string) *Builder {
	b = b.Clone()
	b.joins = b.joins.Join(table, localColumn, operator, foreignColumn)
	return b
}

// LeftJoin adds a left join clause to the query.
func (b *Builder) LeftJoin(table, localColumn, operator, foreignColumn string) *Builder {
	b = b.Clone()
	b.joins = b.joins.LeftJoin(table, localColumn, operator, foreignColumn)
	return b
}

// RightJoin adds a right join clause to the query.
func (b *Builder) RightJoin(table, localColumn, operator, foreignColumn string) *Builder {
	b = b.Clone()
	b.joins = b.joins.RightJoin(table, localColumn, operator, foreignColumn)
	return b
}

// InnerJoin adds an inner join clause to the query.
func (b *Builder) InnerJoin(table, localColumn, operator, foreignColumn string) *Builder {
	b = b.Clone()
	b.joins = b.joins.InnerJoin(table, localColumn, operator, foreignColumn)
	return b
}

// CrossJoin adds a cross join clause to the query.
func (b *Builder) CrossJoin(table, localColumn, operator, foreignColumn string) *Builder {
	b = b.Clone()
	b.joins = b.joins.CrossJoin(table, localColumn, operator, foreignColumn)
	return b
}

// JoinOn adds a join clause to the query with a complex on statement.
func (b *Builder) JoinOn(table string, cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.joins = b.joins.JoinOn(table, cb)
	return b
}

// LeftJoinOn adds a left join clause to the query with a complex on statement.
func (b *Builder) LeftJoinOn(table string, cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.joins = b.joins.LeftJoinOn(table, cb)
	return b
}

// RightJoinOn adds a right join clause to the query with a complex on statement.
func (b *Builder) RightJoinOn(table string, cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.joins = b.joins.RightJoinOn(table, cb)
	return b
}

// InnerJoinOn adds an inner join clause to the query with a complex on statement.
func (b *Builder) InnerJoinOn(table string, cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.joins = b.joins.InnerJoinOn(table, cb)
	return b
}

// CrossJoinOn adds a cross join clause to the query with a complex on statement.
func (b *Builder) CrossJoinOn(table string, cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.joins = b.joins.CrossJoinOn(table, cb)
	return b
}

// Limit set the maximum number of rows to return.
func (b *Builder) Limit(limit int) *Builder {
	b = b.Clone()
	b.limit = b.limit.Limit(limit)
	return b
}

// Offset sets the number of rows to skip before returning the result.
func (b *Builder) Offset(offset int) *Builder {
	b = b.Clone()
	b.limit = b.limit.Offset(offset)
	return b
}

// OrderBy adds an order by clause to the query.
func (b *Builder) OrderBy(column string) *Builder {
	b = b.Clone()
	b.orderBys = b.orderBys.OrderBy(column)
	return b
}

// OrderByDesc adds a descending order by clause to the query.
func (b *Builder) OrderByDesc(column string) *Builder {
	b = b.Clone()
	b.orderBys = b.orderBys.OrderByDesc(column)
	return b
}

// Unordered removes all order by clauses from the query.
func (b *Builder) Unordered() *Builder {
	b = b.Clone()
	b.orderBys = b.orderBys.Unordered()
	return b
}

// WithScope adds a local scope to a query.
func (b *Builder) WithScope(scope *Scope) *Builder {
	b = b.Clone()
	b.scopes = b.scopes.WithScope(scope)
	return b
}

// WithoutScope removes the given scope from the local scopes.
func (b *Builder) WithoutScope(scope *Scope) *Builder {
	b = b.Clone()
	b.scopes = b.scopes.WithoutScope(scope)
	return b
}

// WithoutGlobalScope removes a global scope from the query.
func (b *Builder) WithoutGlobalScope(scope *Scope) *Builder {
	b = b.Clone()
	b.scopes = b.scopes.WithoutGlobalScope(scope)
	return b
}

// Select sets the columns to be selected.
func (b *Builder) Select(columns ...string) *Builder {
	b = b.Clone()
	b.selects = b.selects.Select(columns...)
	return b
}

// AddSelect adds new columns to be selected.
func (b *Builder) AddSelect(columns ...string) *Builder {
	b = b.Clone()
	b.selects = b.selects.AddSelect(columns...)
	return b
}

// SelectSubquery sets a subquery to be selected.
func (b *Builder) SelectSubquery(sb QueryBuilder, as string) *Builder {
	b = b.Clone()
	b.selects = b.selects.SelectSubquery(sb, as)
	return b
}

// AddSelectSubquery adds a subquery to be selected.
func (b *Builder) AddSelectSubquery(sb QueryBuilder, as string) *Builder {
	b = b.Clone()
	b.selects = b.selects.AddSelectSubquery(sb, as)
	return b
}

// SelectFunction sets a column to be selected with a function applied.
func (b *Builder) SelectFunction(function, column string) *Builder {
	b = b.Clone()
	b.selects = b.selects.SelectFunction(function, column)
	return b
}

// SelectFunction adds a column to be selected with a function applied.
func (b *Builder) AddSelectFunction(function, column string) *Builder {
	b = b.Clone()
	b.selects = b.selects.AddSelectFunction(function, column)
	return b
}

// Distinct forces the query to only return distinct results.
func (b *Builder) Distinct() *Builder {
	b = b.Clone()
	b.selects = b.selects.Distinct()
	return b
}

// Where adds a basic where clause to the query.
func (b *Builder) Where(column, operator string, value any) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.Where(column, operator, value)
	return b
}

// Having adds a basic having clause to the query.
func (b *Builder) Having(column, operator string, value any) *Builder {
	b = b.Clone()
	b.havings = b.havings.Where(column, operator, value)
	return b
}

// OrWhere adds an or where clause to the query
func (b *Builder) OrWhere(column, operator string, value any) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhere(column, operator, value)
	return b
}

// OrHaving adds an or having clause to the query
func (b *Builder) OrHaving(column, operator string, value any) *Builder {
	b = b.Clone()
	b.havings = b.havings.OrWhere(column, operator, value)
	return b
}

// WhereColumn adds a where clause to the query comparing two columns.
func (b *Builder) WhereColumn(column, operator string, valueColumn string) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.WhereColumn(column, operator, valueColumn)
	return b
}

// HavingColumn adds a having clause to the query comparing two columns.
func (b *Builder) HavingColumn(column, operator string, valueColumn string) *Builder {
	b = b.Clone()
	b.havings = b.havings.WhereColumn(column, operator, valueColumn)
	return b
}

// OrWhereColumn adds an or where clause to the query comparing two columns.
func (b *Builder) OrWhereColumn(column, operator string, valueColumn string) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereColumn(column, operator, valueColumn)
	return b
}

// OrHavingColumn adds an or having clause to the query comparing two columns.
func (b *Builder) OrHavingColumn(column, operator string, valueColumn string) *Builder {
	b = b.Clone()
	b.havings = b.havings.OrWhereColumn(column, operator, valueColumn)
	return b
}

// WhereIn adds a where in clause to the query.
func (b *Builder) WhereIn(column string, values []any) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.WhereIn(column, values)
	return b
}

// HavingIn adds a having in clause to the query.
func (b *Builder) HavingIn(column string, values []any) *Builder {
	b = b.Clone()
	b.havings = b.havings.WhereIn(column, values)
	return b
}

// OrWhereIn adds an or where in clause to the query.
func (b *Builder) OrWhereIn(column string, values []any) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereIn(column, values)
	return b
}

// OrHavingIn adds an or having in clause to the query.
func (b *Builder) OrHavingIn(column string, values []any) *Builder {
	b = b.Clone()
	b.havings = b.havings.OrWhereIn(column, values)
	return b
}

// WhereExists add an exists clause to the query.
func (b *Builder) WhereExists(query QueryBuilder) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.WhereExists(query)
	return b
}

// HavingExists add an exists clause to the query.
func (b *Builder) HavingExists(query QueryBuilder) *Builder {
	b = b.Clone()
	b.havings = b.havings.WhereExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *Builder) OrWhereExists(query QueryBuilder) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *Builder) OrHavingExists(query QueryBuilder) *Builder {
	b = b.Clone()
	b.havings = b.havings.OrWhereExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *Builder) WhereNotExists(query QueryBuilder) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.WhereNotExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *Builder) HavingNotExists(query QueryBuilder) *Builder {
	b = b.Clone()
	b.havings = b.havings.WhereNotExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *Builder) OrWhereNotExists(query QueryBuilder) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereNotExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *Builder) OrHavingNotExists(query QueryBuilder) *Builder {
	b = b.Clone()
	b.havings = b.havings.OrWhereNotExists(query)
	return b
}

// WhereSubquery adds a where clause to the query comparing a column and a subquery.
func (b *Builder) WhereSubquery(subquery QueryBuilder, operator string, value any) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.WhereSubquery(subquery, operator, value)
	return b
}

// HavingSubquery adds a having clause to the query comparing a column and a subquery.
func (b *Builder) HavingSubquery(subquery QueryBuilder, operator string, value any) *Builder {
	b = b.Clone()
	b.havings = b.havings.WhereSubquery(subquery, operator, value)
	return b
}

// OrWhereSubquery adds an or where clause to the query comparing a column and a subquery.
func (b *Builder) OrWhereSubquery(subquery QueryBuilder, operator string, value any) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereSubquery(subquery, operator, value)
	return b
}

// OrHavingSubquery adds an or having clause to the query comparing a column and a subquery.
func (b *Builder) OrHavingSubquery(subquery QueryBuilder, operator string, value any) *Builder {
	b = b.Clone()
	b.havings = b.havings.OrWhereSubquery(subquery, operator, value)
	return b
}

// WhereHas adds a relationship exists condition to the query with where clauses.
func (b *Builder) WhereHas(relation string, cb func(q *Builder) *Builder) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.WhereHas(relation, cb)
	return b
}

// HavingHas adds a relationship exists condition to the query with having clauses.
func (b *Builder) HavingHas(relation string, cb func(q *Builder) *Builder) *Builder {
	b = b.Clone()
	b.havings = b.havings.WhereHas(relation, cb)
	return b
}

// OrWhereHas adds a relationship exists condition to the query with where clauses and an or.
func (b *Builder) OrWhereHas(relation string, cb func(q *Builder) *Builder) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereHas(relation, cb)
	return b
}

// OrHavingHas adds a relationship exists condition to the query with having clauses and an or.
func (b *Builder) OrHavingHas(relation string, cb func(q *Builder) *Builder) *Builder {
	b = b.Clone()
	b.havings = b.havings.OrWhereHas(relation, cb)
	return b
}

// WhereRaw adds a raw where clause to the query.
func (b *Builder) WhereRaw(rawSql string, bindings ...any) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.WhereRaw(rawSql, bindings...)
	return b
}

// HavingRaw adds a raw having clause to the query.
func (b *Builder) HavingRaw(rawSql string, bindings ...any) *Builder {
	b = b.Clone()
	b.havings = b.havings.WhereRaw(rawSql, bindings...)
	return b
}

// OrWhereRaw adds a raw or where clause to the query.
func (b *Builder) OrWhereRaw(rawSql string, bindings ...any) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereRaw(rawSql, bindings...)
	return b
}

// OrHavingRaw adds a raw or having clause to the query.
func (b *Builder) OrHavingRaw(rawSql string, bindings ...any) *Builder {
	b = b.Clone()
	b.havings = b.havings.OrWhereRaw(rawSql, bindings...)
	return b
}

// And adds a group of conditions to the query
func (b *Builder) And(cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.And(cb)
	return b
}

// HavingAnd adds a group of conditions to the query
func (b *Builder) HavingAnd(cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.havings = b.havings.And(cb)
	return b
}

// Or adds a group of conditions to the query with an or
func (b *Builder) Or(cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.wheres = b.wheres.Or(cb)
	return b
}

// HavingOr adds a group of conditions to the query with an or
func (b *Builder) HavingOr(cb func(q *Conditions)) *Builder {
	b = b.Clone()
	b.havings = b.havings.Or(cb)
	return b
}
