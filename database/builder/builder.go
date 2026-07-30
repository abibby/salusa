package builder

import (
	"context"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/extra/sets"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/internal/relationship"
)

// QueryBuilder is implemented by *ModelBuilder and *Builder
type QueryBuilder interface {
	helpers.SQLStringer
	imALittleQueryBuilderShortAndStout()
}

type Builder struct {
	query dialects.Query
	// selects  *selects
	// from     fromTable
	// joins    joins
	// wheres   *Conditions
	// groupBys groupBys
	// havings  *Conditions
	// limit    *limit
	// orderBys orderBys
	scopes *scopes
	ctx    context.Context
}

var _ dialects.QueryBuilder = (*Builder)(nil)

// NewBuilder creates a new SubBuilder without anything selected
func NewBuilder() *Builder {
	return &Builder{
		query:  dialects.NewQuery(),
		scopes: newScopes(),
		ctx:    context.Background(),
	}
}

// Query implements [dialects.QueryBuilder].
func (b *Builder) Query() *dialects.Query {
	return &b.query
}

// ModelBuilder represents an sql query and any bindings needed to run it.
//
//go:generate go run ../../internal/build/build.go
type ModelBuilder[T model.Model] struct {
	builder       *Builder
	withs         []string
	withoutScopes sets.Set[string]
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

	_ = relationship.InitializeRelationships(m)

	sb := NewBuilder()
	sb.scopes.withParent(m)
	return &ModelBuilder[T]{
		builder:       sb,
		withs:         []string{},
		withoutScopes: sets.New[string](),
	}
}

// Query implements [dialects.QueryBuilder].
func (b *ModelBuilder[T]) Query() *dialects.Query {
	return b.builder.Query()
}

func (*ModelBuilder[T]) imALittleQueryBuilderShortAndStout() {}
func (*Builder) imALittleQueryBuilderShortAndStout()         {}
