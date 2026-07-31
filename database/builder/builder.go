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

//go:generate go run ../../internal/build/build.go
type Builder struct {
	query   dialects.SelectQuery
	wheres  *Conditions
	havings *Conditions

	scopes *scopes
	ctx    context.Context
}

var _ dialects.QueryBuilder = (*Builder)(nil)

// NewBuilder creates a new SubBuilder without anything selected
func NewBuilder() *Builder {
	return &Builder{
		query:   dialects.NewSelectQuery(),
		wheres:  newConditions(),
		havings: newConditions(),
		scopes:  newScopes(),
		ctx:     context.Background(),
	}
}

// Query implements [dialects.QueryBuilder].
func (b *Builder) Query() *dialects.SelectQuery {
	q := &b.query
	q.Wheres = b.wheres.conditions
	q.Havings = b.havings.conditions
	return q
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
	sb.wheres.withParent(m)
	sb.havings.withParent(m)
	sb.scopes.withParent(m)

	return &ModelBuilder[T]{
		builder:       sb,
		withs:         []string{},
		withoutScopes: sets.New[string](),
	}
}

// Query implements [dialects.QueryBuilder].
func (b *ModelBuilder[T]) Query() *dialects.SelectQuery {
	return b.builder.Query()
}

func (*ModelBuilder[T]) imALittleQueryBuilderShortAndStout() {}
func (*Builder) imALittleQueryBuilderShortAndStout()         {}
