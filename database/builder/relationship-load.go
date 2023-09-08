package builder

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/internal/relationship"
)

func Load(tx helpers.QueryExecer, v any, relation string) error {
	return LoadContext(context.Background(), tx, v, relation)
}
func LoadContext(ctx context.Context, tx helpers.QueryExecer, v any, relation string) error {
	return loadContext(ctx, tx, v, relation, false)
}
func LoadMissing(tx helpers.QueryExecer, v any, relation string) error {
	return LoadMissingContext(context.Background(), tx, v, relation)
}
func LoadMissingContext(ctx context.Context, tx helpers.QueryExecer, v any, relation string) error {
	return loadContext(ctx, tx, v, relation, true)
}
func loadContext(ctx context.Context, tx helpers.QueryExecer, v any, relation string, onlyMissing bool) error {
	err := relationship.InitializeRelationships(v)
	if err != nil {
		return err
	}
	relations := strings.Split(relation, ".")
	for i, rel := range relations {
		relations := []Relationship{}
		err := helpers.Each(v, func(v reflect.Value, pointer bool) error {
			r, ok := getRelation(v, rel)
			if !ok {
				return fmt.Errorf("%s has no relation %s: %w", v.Type().Name(), rel, ErrMissingRelationship)
			}

			if onlyMissing && r.Loaded() {
				return nil
			}
			relations = append(relations, r)
			return nil
		})
		if err != nil {
			return err
		}

		if len(relations) == 0 {
			return nil
		}

		err = relations[0].Load(ctx, tx, relations)
		if err != nil {
			return err
		}

		if i <= len(relations)-1 {
			values := []any{}
			err := helpers.Each(v, func(v reflect.Value, pointer bool) error {
				related, ok := getValue(v, rel)
				if !ok {
					return nil
				}
				out := related.MethodByName("Value").Call([]reflect.Value{})
				if len(out) != 2 {
					return fmt.Errorf("invalid value method")
				}
				if !out[1].IsZero() {
					values = append(values, out[0].Interface())
				}
				return nil
			})
			if err != nil {
				return err
			}
			v = values
		}
	}

	return nil
}

func getValue(rv reflect.Value, key string) (reflect.Value, bool) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if ft.Anonymous {
			result, ok := getValue(rv.Field(i), key)
			if ok {
				return result, true
			}
			continue
		}
		if ft.Name == key {
			return rv.Field(i), true
		}
	}
	return reflect.Value{}, false
}
func ofType[T Relationship](relations []Relationship) []T {
	relationsOfType := make([]T, 0, len(relations))
	for _, r := range relations {
		rOfType, ok := r.(T)
		if ok {
			relationsOfType = append(relationsOfType, rOfType)
		}
	}
	return relationsOfType
}
