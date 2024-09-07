package builder

import (
	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/set"
)

type Scoper interface {
	Scopes() []*Scope
}

// Scope is a modifier for a query that can be easily applied.
type Scope struct {
	Name   string
	Query  ScopeQueryFunc
	Delete ScopeDeleteFunc
}
type ScopeQueryFunc func(b *Builder) *Builder
type ScopeDeleteFunc func(next func(q *Builder, tx database.DB) error) func(q *Builder, tx database.DB) error

type scopes struct {
	parent              any
	scopes              []*Scope
	withoutGlobalScopes set.Set[string]
}

func newScopes() *scopes {
	return &scopes{
		scopes:              []*Scope{},
		withoutGlobalScopes: set.New[string](),
	}
}

func (s *scopes) withParent(parnet any) *scopes {
	s.parent = parnet
	return s
}

func (s *scopes) Clone() *scopes {
	return &scopes{
		parent:              s.parent,
		scopes:              cloneSlice(s.scopes),
		withoutGlobalScopes: s.withoutGlobalScopes.Clone(),
	}
}

// WithScope adds a local scope to a Delete.
func (s *scopes) WithScope(scope *Scope) *scopes {
	s.scopes = append(s.scopes, scope)
	return s
}

// WithoutScope removes the given scope from the local scopes.
func (s *scopes) WithoutScope(scope *Scope) *scopes {
	newScopes := make([]*Scope, 0, len(s.scopes))
	for _, sc := range s.scopes {
		if sc.Name != scope.Name {
			newScopes = append(newScopes, sc)
		}
	}
	s.scopes = newScopes
	return s
}

func (s *scopes) allScopes() []*Scope {
	if scoper, ok := s.parent.(Scoper); ok {
		allGlobalScopes := scoper.Scopes()
		combinedScopes := make([]*Scope, len(s.scopes), len(s.scopes)+len(allGlobalScopes))
		copy(combinedScopes, s.scopes)

		for _, scope := range allGlobalScopes {
			if s.withoutGlobalScopes.Has(scope.Name) {
				continue
			}
			combinedScopes = append(combinedScopes, scope)
		}
		return combinedScopes
	}

	return s.scopes
}

// WithoutGlobalScope removes a global scope from the query.
func (b *scopes) WithoutGlobalScope(scope *Scope) *scopes {
	b.withoutGlobalScopes.Add(scope.Name)
	return b
}

func (b *Builder) ActiveScopes() []*Scope {
	return b.scopes.allScopes()
}
func (b *ModelBuilder[T]) ActiveScopes() []*Scope {
	return b.builder.ActiveScopes()
}
