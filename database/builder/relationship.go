package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/internal/relationship"
)

type ForeignKey struct {
	LocalKey     string
	RelatedTable string
	RelatedKey   string
}

func (f *ForeignKey) Equal(v *ForeignKey) bool {
	return f.LocalKey == v.LocalKey && f.RelatedKey == v.RelatedKey && f.RelatedTable == v.RelatedTable
}

type Relationship interface {
	relationship.Relationship
	Subquery() *SubBuilder
	Load(ctx context.Context, tx helpers.QueryExecer, relations []Relationship) error
	ForeignKeys() []*ForeignKey
}

type relationValue[T any] struct {
	loaded bool
	value  T
}

var (
	ErrMissingRelationship = fmt.Errorf("missing relationship")
	ErrMissingField        = fmt.Errorf("missing related field")
)

// Value will return the related value and if it has been fetched.
func (v *relationValue[T]) Value() (T, bool) {
	return v.value, v.loaded
}

// Loaded returns true if the relationship has been fetched and false if it has
// not.
func (v *relationValue[T]) Loaded() bool {
	return v.loaded
}

func (v *relationValue[T]) MarshalJSON() ([]byte, error) {
	if !v.loaded {
		return json.Marshal(nil)
	}
	return json.Marshal(v.value)
}

func foreignKeyName(field reflect.StructField, tag string, tableType any) (string, error) {
	t, ok := field.Tag.Lookup(tag)
	if ok {
		return t, nil
	}

	pKeys := helpers.PrimaryKey(tableType)
	if len(pKeys) != 1 {
		return "", fmt.Errorf("you must specify keys for relationships with compound primary keys")
	}
	return database.GetTableSingular(tableType) + "_" + pKeys[0], nil
}

func primaryKeyName(field reflect.StructField, tag string, tableType any) (string, error) {
	t, ok := field.Tag.Lookup(tag)
	if ok {
		return t, nil
	}
	pKeys := helpers.PrimaryKey(tableType)
	if len(pKeys) != 1 {
		return "", fmt.Errorf("you must specify keys for relationships with compound primary keys")
	}
	return pKeys[0], nil
}

func getRelation(rv reflect.Value, relation string) (Relationship, bool) {
	if rv.Kind() == reflect.Ptr {
		if rv.IsZero() {
			rv = reflect.New(rv.Type().Elem())
			err := relationship.InitializeRelationships(rv.Interface())
			if err != nil {
				panic(err)
			}
		}
		rv = rv.Elem()
	}
	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		ft := t.Field(i)

		if ft.Anonymous {
			r, ok := getRelation(rv.Field(i), relation)
			if ok {
				return r, true
			}
			continue
		}

		if ft.Name != relation {
			continue
		}

		r, ok := rv.Field(i).Interface().(Relationship)
		if !ok {
			continue
		}
		return r, true
	}
	return nil, false
}
