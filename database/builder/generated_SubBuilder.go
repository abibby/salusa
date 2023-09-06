package builder

// From sets the table which the query is targeting.
func (b *SubBuilder) From(table string) *SubBuilder {
	b = b.Clone()
	b.from = b.from.From(table)
	return b
}

// GroupBy sets the "group by" clause to the query.
func (b *SubBuilder) GroupBy(columns ...string) *SubBuilder {
	b = b.Clone()
	b.groupBys = b.groupBys.GroupBy(columns...)
	return b
}

// GroupBy adds a "group by" clause to the query.
func (b *SubBuilder) AddGroupBy(columns ...string) *SubBuilder {
	b = b.Clone()
	b.groupBys = b.groupBys.AddGroupBy(columns...)
	return b
}

// Join adds a join clause to the query.
func (b *SubBuilder) Join(table, localColumn, operator, foreignColumn string) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.Join(table, localColumn, operator, foreignColumn)
	return b
}

// LeftJoin adds a left join clause to the query.
func (b *SubBuilder) LeftJoin(table, localColumn, operator, foreignColumn string) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.LeftJoin(table, localColumn, operator, foreignColumn)
	return b
}

// RightJoin adds a right join clause to the query.
func (b *SubBuilder) RightJoin(table, localColumn, operator, foreignColumn string) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.RightJoin(table, localColumn, operator, foreignColumn)
	return b
}

// InnerJoin adds an inner join clause to the query.
func (b *SubBuilder) InnerJoin(table, localColumn, operator, foreignColumn string) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.InnerJoin(table, localColumn, operator, foreignColumn)
	return b
}

// CrossJoin adds a cross join clause to the query.
func (b *SubBuilder) CrossJoin(table, localColumn, operator, foreignColumn string) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.CrossJoin(table, localColumn, operator, foreignColumn)
	return b
}

// JoinOn adds a join clause to the query with a complex on statement.
func (b *SubBuilder) JoinOn(table string, cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.JoinOn(table, cb)
	return b
}

// LeftJoinOn adds a left join clause to the query with a complex on statement.
func (b *SubBuilder) LeftJoinOn(table string, cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.LeftJoinOn(table, cb)
	return b
}

// RightJoinOn adds a right join clause to the query with a complex on statement.
func (b *SubBuilder) RightJoinOn(table string, cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.RightJoinOn(table, cb)
	return b
}

// InnerJoinOn adds an inner join clause to the query with a complex on statement.
func (b *SubBuilder) InnerJoinOn(table string, cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.InnerJoinOn(table, cb)
	return b
}

// CrossJoinOn adds a cross join clause to the query with a complex on statement.
func (b *SubBuilder) CrossJoinOn(table string, cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.joins = b.joins.CrossJoinOn(table, cb)
	return b
}

// Limit set the maximum number of rows to return.
func (b *SubBuilder) Limit(limit int) *SubBuilder {
	b = b.Clone()
	b.limit = b.limit.Limit(limit)
	return b
}

// Offset sets the number of rows to skip before returning the result.
func (b *SubBuilder) Offset(offset int) *SubBuilder {
	b = b.Clone()
	b.limit = b.limit.Offset(offset)
	return b
}

// OrderBy adds an order by clause to the query.
func (b *SubBuilder) OrderBy(column string) *SubBuilder {
	b = b.Clone()
	b.orderBys = b.orderBys.OrderBy(column)
	return b
}

// OrderByDesc adds a descending order by clause to the query.
func (b *SubBuilder) OrderByDesc(column string) *SubBuilder {
	b = b.Clone()
	b.orderBys = b.orderBys.OrderByDesc(column)
	return b
}

// Unordered removes all order by clauses from the query.
func (b *SubBuilder) Unordered() *SubBuilder {
	b = b.Clone()
	b.orderBys = b.orderBys.Unordered()
	return b
}

// WithScope adds a local scope to a query.
func (b *SubBuilder) WithScope(scope *Scope) *SubBuilder {
	b = b.Clone()
	b.scopes = b.scopes.WithScope(scope)
	return b
}

// WithoutScope removes the given scope from the local scopes.
func (b *SubBuilder) WithoutScope(scope *Scope) *SubBuilder {
	b = b.Clone()
	b.scopes = b.scopes.WithoutScope(scope)
	return b
}

// WithoutGlobalScope removes a global scope from the query.
func (b *SubBuilder) WithoutGlobalScope(scope *Scope) *SubBuilder {
	b = b.Clone()
	b.scopes = b.scopes.WithoutGlobalScope(scope)
	return b
}

// Select sets the columns to be selected.
func (b *SubBuilder) Select(columns ...string) *SubBuilder {
	b = b.Clone()
	b.selects = b.selects.Select(columns...)
	return b
}

// AddSelect adds new columns to be selected.
func (b *SubBuilder) AddSelect(columns ...string) *SubBuilder {
	b = b.Clone()
	b.selects = b.selects.AddSelect(columns...)
	return b
}

// SelectSubquery sets a subquery to be selected.
func (b *SubBuilder) SelectSubquery(sb QueryBuilder, as string) *SubBuilder {
	b = b.Clone()
	b.selects = b.selects.SelectSubquery(sb, as)
	return b
}

// AddSelectSubquery adds a subquery to be selected.
func (b *SubBuilder) AddSelectSubquery(sb QueryBuilder, as string) *SubBuilder {
	b = b.Clone()
	b.selects = b.selects.AddSelectSubquery(sb, as)
	return b
}

// SelectFunction sets a column to be selected with a function applied.
func (b *SubBuilder) SelectFunction(function, column string) *SubBuilder {
	b = b.Clone()
	b.selects = b.selects.SelectFunction(function, column)
	return b
}

// SelectFunction adds a column to be selected with a function applied.
func (b *SubBuilder) AddSelectFunction(function, column string) *SubBuilder {
	b = b.Clone()
	b.selects = b.selects.AddSelectFunction(function, column)
	return b
}

// Distinct forces the query to only return distinct results.
func (b *SubBuilder) Distinct() *SubBuilder {
	b = b.Clone()
	b.selects = b.selects.Distinct()
	return b
}

// Where adds a basic where clause to the query.
func (b *SubBuilder) Where(column, operator string, value any) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.Where(column, operator, value)
	return b
}

// Having adds a basic having clause to the query.
func (b *SubBuilder) Having(column, operator string, value any) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.Where(column, operator, value)
	return b
}

// OrWhere adds an or where clause to the query
func (b *SubBuilder) OrWhere(column, operator string, value any) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhere(column, operator, value)
	return b
}

// OrHaving adds an or having clause to the query
func (b *SubBuilder) OrHaving(column, operator string, value any) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.OrWhere(column, operator, value)
	return b
}

// WhereColumn adds a where clause to the query comparing two columns.
func (b *SubBuilder) WhereColumn(column, operator string, valueColumn string) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.WhereColumn(column, operator, valueColumn)
	return b
}

// HavingColumn adds a having clause to the query comparing two columns.
func (b *SubBuilder) HavingColumn(column, operator string, valueColumn string) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.WhereColumn(column, operator, valueColumn)
	return b
}

// OrWhereColumn adds an or where clause to the query comparing two columns.
func (b *SubBuilder) OrWhereColumn(column, operator string, valueColumn string) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereColumn(column, operator, valueColumn)
	return b
}

// OrHavingColumn adds an or having clause to the query comparing two columns.
func (b *SubBuilder) OrHavingColumn(column, operator string, valueColumn string) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.OrWhereColumn(column, operator, valueColumn)
	return b
}

// WhereIn adds a where in clause to the query.
func (b *SubBuilder) WhereIn(column string, values []any) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.WhereIn(column, values)
	return b
}

// HavingIn adds a having in clause to the query.
func (b *SubBuilder) HavingIn(column string, values []any) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.WhereIn(column, values)
	return b
}

// OrWhereIn adds an or where in clause to the query.
func (b *SubBuilder) OrWhereIn(column string, values []any) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereIn(column, values)
	return b
}

// OrHavingIn adds an or having in clause to the query.
func (b *SubBuilder) OrHavingIn(column string, values []any) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.OrWhereIn(column, values)
	return b
}

// WhereExists add an exists clause to the query.
func (b *SubBuilder) WhereExists(query QueryBuilder) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.WhereExists(query)
	return b
}

// HavingExists add an exists clause to the query.
func (b *SubBuilder) HavingExists(query QueryBuilder) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.WhereExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *SubBuilder) OrWhereExists(query QueryBuilder) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereExists(query)
	return b
}

// WhereExists add an exists clause to the query.
func (b *SubBuilder) OrHavingExists(query QueryBuilder) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.OrWhereExists(query)
	return b
}

// WhereSubquery adds a where clause to the query comparing a column and a subquery.
func (b *SubBuilder) WhereSubquery(subquery QueryBuilder, operator string, value any) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.WhereSubquery(subquery, operator, value)
	return b
}

// HavingSubquery adds a having clause to the query comparing a column and a subquery.
func (b *SubBuilder) HavingSubquery(subquery QueryBuilder, operator string, value any) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.WhereSubquery(subquery, operator, value)
	return b
}

// OrWhereSubquery adds an or where clause to the query comparing a column and a subquery.
func (b *SubBuilder) OrWhereSubquery(subquery QueryBuilder, operator string, value any) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereSubquery(subquery, operator, value)
	return b
}

// OrHavingSubquery adds an or having clause to the query comparing a column and a subquery.
func (b *SubBuilder) OrHavingSubquery(subquery QueryBuilder, operator string, value any) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.OrWhereSubquery(subquery, operator, value)
	return b
}

// WhereHas adds a relationship exists condition to the query with where clauses.
func (b *SubBuilder) WhereHas(relation string, cb func(q *SubBuilder) *SubBuilder) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.WhereHas(relation, cb)
	return b
}

// HavingHas adds a relationship exists condition to the query with having clauses.
func (b *SubBuilder) HavingHas(relation string, cb func(q *SubBuilder) *SubBuilder) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.WhereHas(relation, cb)
	return b
}

// OrWhereHas adds a relationship exists condition to the query with where clauses and an or.
func (b *SubBuilder) OrWhereHas(relation string, cb func(q *SubBuilder) *SubBuilder) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereHas(relation, cb)
	return b
}

// OrHavingHas adds a relationship exists condition to the query with having clauses and an or.
func (b *SubBuilder) OrHavingHas(relation string, cb func(q *SubBuilder) *SubBuilder) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.OrWhereHas(relation, cb)
	return b
}

// WhereRaw adds a raw where clause to the query.
func (b *SubBuilder) WhereRaw(rawSql string, bindings ...any) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.WhereRaw(rawSql, bindings...)
	return b
}

// HavingRaw adds a raw having clause to the query.
func (b *SubBuilder) HavingRaw(rawSql string, bindings ...any) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.WhereRaw(rawSql, bindings...)
	return b
}

// OrWhereRaw adds a raw or where clause to the query.
func (b *SubBuilder) OrWhereRaw(rawSql string, bindings ...any) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.OrWhereRaw(rawSql, bindings...)
	return b
}

// OrHavingRaw adds a raw or having clause to the query.
func (b *SubBuilder) OrHavingRaw(rawSql string, bindings ...any) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.OrWhereRaw(rawSql, bindings...)
	return b
}

// And adds a group of conditions to the query
func (b *SubBuilder) And(cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.And(cb)
	return b
}

// HavingAnd adds a group of conditions to the query
func (b *SubBuilder) HavingAnd(cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.And(cb)
	return b
}

// Or adds a group of conditions to the query with an or
func (b *SubBuilder) Or(cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.wheres = b.wheres.Or(cb)
	return b
}

// HavingOr adds a group of conditions to the query with an or
func (b *SubBuilder) HavingOr(cb func(q *Conditions)) *SubBuilder {
	b = b.Clone()
	b.havings = b.havings.Or(cb)
	return b
}
