package builder

import (
	"context"

	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/set"
)

// QueryBuilder is implemented by *Builder and *SubBuilder
type QueryBuilder interface {
	helpers.ToSQLer
	imALittleQueryBuilderShortAndStout()
}

//go:generate go run ../../internal/build/build.go
type SubBuilder struct {
	selects  *selects
	from     fromTable
	joins    joins
	wheres   *Conditions
	groupBys groupBys
	havings  *Conditions
	limit    *limit
	orderBys orderBys
	scopes   *scopes
	ctx      context.Context
}

// Builder represents an sql query and any bindings needed to run it.
//
//go:generate go run ../../internal/build/build.go
type Builder[T model.Model] struct {
	subBuilder    *SubBuilder
	withs         []string
	withoutScopes set.Set[string]
}

// New creates a new Builder with * selected
func New[T model.Model]() *Builder[T] {
	return NewEmpty[T]().Select("*")
}

// From creates a new query from the models table and with table.* selected
func From[T model.Model]() *Builder[T] {
	var m T
	table := helpers.GetTable(m)
	return NewEmpty[T]().Select(table + ".*").From(table)
}

// NewEmpty creates a new helpers without anything selected
func NewEmpty[T model.Model]() *Builder[T] {
	var m T
	sb := NewSubBuilder()
	sb.wheres.withParent(m)
	sb.havings.withParent(m)
	sb.scopes.withParent(m)
	return &Builder[T]{
		subBuilder:    sb,
		withs:         []string{},
		withoutScopes: set.New[string](),
	}
}

// NewSubBuilder creates a new SubBuilder without anything selected
func NewSubBuilder() *SubBuilder {
	return &SubBuilder{
		selects:  NewSelects(),
		from:     "",
		wheres:   newConditions().withPrefix("WHERE"),
		groupBys: groupBys{},
		havings:  newConditions().withPrefix("HAVING"),
		limit:    &limit{},
		scopes:   newScopes(),
		ctx:      context.Background(),
	}
}

func (*Builder[T]) imALittleQueryBuilderShortAndStout() {}
func (*SubBuilder) imALittleQueryBuilderShortAndStout() {}
