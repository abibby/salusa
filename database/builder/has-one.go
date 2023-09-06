package builder

import (
	"context"
	"reflect"

	"github.com/abibby/salusa/database/models"
	"github.com/abibby/salusa/internal/helpers"
)

// # Tags:
//   - local: parent model
//   - foreign: related model
type HasOne[T models.Model] struct {
	hasOneOrMany[T]
	relationValue[T]
}

var _ Relationship = &HasOne[models.Model]{}

func (r *HasOne[T]) Initialize(parent any, field reflect.StructField) error {
	r.parent = parent
	parentKey, err := primaryKeyName(field, "local", parent)
	if err != nil {
		return err
	}
	relatedKey, err := foreignKeyName(field, "foreign", parent)
	if err != nil {
		return err
	}
	r.parentKey = parentKey
	r.relatedKey = relatedKey
	return nil
}

func (r *HasOne[T]) Load(ctx context.Context, tx helpers.QueryExecer, relations []Relationship) error {
	rm, err := r.relatedMap(ctx, tx, relations)
	if err != nil {
		return err
	}

	for _, relation := range ofType[*HasOne[T]](relations) {
		relation.value = rm.Single(relation.parentKeyValue())
		relation.loaded = true
	}
	return nil
}

// ForeignKeys returns a list of related tables and what columns they are
// related on.
func (r *HasOne[T]) ForeignKeys() []*ForeignKey {
	return []*ForeignKey{}
}
