package builder

import (
	"context"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/set"
)

// QueryBuilder is implemented by *ModelBuilder and *Builder
type QueryBuilder interface {
	helpers.SQLStringer
	imALittleQueryBuilderShortAndStout()
}

//go:generate go run ../../internal/build/build.go
type Builder struct {
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

// ModelBuilder represents an sql query and any bindings needed to run it.
//
//go:generate go run ../../internal/build/build.go
type ModelBuilder[T model.Model] struct {
	builder       *Builder
	withs         []string
	withoutScopes set.Set[string]
}

// New creates a new Builder with * selected
func New[T model.Model]() *ModelBuilder[T] {
	return NewEmpty[T]().Select("*")
}

// From creates a new query from the models table and with table.* selected
func From[T model.Model]() *ModelBuilder[T] {
	var m T
	table := database.GetTable(m)
	return NewEmpty[T]().Select(table + ".*").From(table)
}

// NewEmpty creates a new helpers without anything selected
func NewEmpty[T model.Model]() *ModelBuilder[T] {
	m := helpers.CreateFor[T]().Interface().(T)
	sb := NewBuilder()
	sb.wheres.withParent(m)
	sb.havings.withParent(m)
	sb.scopes.withParent(m)
	return &ModelBuilder[T]{
		builder:       sb,
		withs:         []string{},
		withoutScopes: set.New[string](),
	}
}

// NewBuilder creates a new SubBuilder without anything selected
func NewBuilder() *Builder {
	return &Builder{
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

func (*ModelBuilder[T]) imALittleQueryBuilderShortAndStout() {}
func (*Builder) imALittleQueryBuilderShortAndStout()         {}
