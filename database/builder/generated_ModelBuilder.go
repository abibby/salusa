package builder

import (
	"context"

	"github.com/abibby/salusa/database/dialects"
)

// WithContext adds a context to the query that will be used when fetching results.
func (b *ModelBuilder[T]) WithContext(ctx context.Context) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.WithContext(ctx)
	return b
}

// From sets the table which the query is targeting.
func (b *ModelBuilder[T]) From(table string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.From(table)
	return b
}

// GroupBy sets the "group by" clause to the query.
func (b *ModelBuilder[T]) GroupBy(columns ...string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.GroupBy(columns...)
	return b
}

// GroupBy adds a "group by" clause to the query.
func (b *ModelBuilder[T]) AddGroupBy(columns ...string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.AddGroupBy(columns...)
	return b
}

// Join adds a join clause to the query.
func (b *ModelBuilder[T]) Join(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.Join(table, localColumn, operator, foreignColumn)
	return b
}

// LeftJoin adds a left join clause to the query.
func (b *ModelBuilder[T]) LeftJoin(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.LeftJoin(table, localColumn, operator, foreignColumn)
	return b
}

// RightJoin adds a right join clause to the query.
func (b *ModelBuilder[T]) RightJoin(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.RightJoin(table, localColumn, operator, foreignColumn)
	return b
}

// InnerJoin adds an inner join clause to the query.
func (b *ModelBuilder[T]) InnerJoin(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.InnerJoin(table, localColumn, operator, foreignColumn)
	return b
}

// CrossJoin adds a cross join clause to the query.
func (b *ModelBuilder[T]) CrossJoin(table, localColumn, operator, foreignColumn string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.CrossJoin(table, localColumn, operator, foreignColumn)
	return b
}

// JoinOn adds a join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) JoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.JoinOn(table, cb)
	return b
}

// LeftJoinOn adds a left join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) LeftJoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.LeftJoinOn(table, cb)
	return b
}

// RightJoinOn adds a right join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) RightJoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.RightJoinOn(table, cb)
	return b
}

// InnerJoinOn adds an inner join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) InnerJoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.InnerJoinOn(table, cb)
	return b
}

// CrossJoinOn adds a cross join clause to the query with a complex on statement.
func (b *ModelBuilder[T]) CrossJoinOn(table string, cb func(q *Conditions)) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.CrossJoinOn(table, cb)
	return b
}

// Select sets the columns to be selected.
func (b *ModelBuilder[T]) Select(columns ...string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.Select(columns...)
	return b
}

// AddSelect adds new columns to be selected.
func (b *ModelBuilder[T]) AddSelect(columns ...string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.AddSelect(columns...)
	return b
}

// SelectSubquery sets a subquery to be selected.
func (b *ModelBuilder[T]) SelectSubquery(sb dialects.QueryBuilder, as string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.SelectSubquery(sb, as)
	return b
}

// AddSelectSubquery adds a subquery to be selected.
func (b *ModelBuilder[T]) AddSelectSubquery(sb dialects.QueryBuilder, as string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.AddSelectSubquery(sb, as)
	return b
}

// SelectFunction sets a column to be selected with a function applied.
func (b *ModelBuilder[T]) SelectFunction(function, column string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.SelectFunction(function, column)
	return b
}

// SelectFunction adds a column to be selected with a function applied.
func (b *ModelBuilder[T]) AddSelectFunction(function, column string) *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.AddSelectFunction(function, column)
	return b
}

// Distinct forces the query to only return distinct results.
func (b *ModelBuilder[T]) Distinct() *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.Distinct()
	return b
}
func (b *ModelBuilder[T]) Dump() *ModelBuilder[T] {
	b = b.Clone()
	b.builder = b.builder.Dump()
	return b
}
