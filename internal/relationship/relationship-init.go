package relationship

import (
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
)

var RelationType = reflect.TypeOf((*Relationship)(nil)).Elem()

func InitializeRelationships(v any) error {
	return helpers.Each(v, initializeRelationships)
}

func initializeRelationships(v reflect.Value, pointer bool) error {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)

		if ft.Anonymous {
			err := initializeRelationships(v.Field(i), ft.Type.Kind() == reflect.Ptr)
			if err != nil {
				return err
			}
			continue
		}

		if ft.Type.Implements(RelationType) {
			fv := v.Field(i)
			if ft.Type.Kind() == reflect.Ptr {
				fv.Set(reflect.New(ft.Type.Elem()))
			} else {
				fv.Set(reflect.New(ft.Type).Elem())
			}
			iValue := v
			if pointer {
				iValue = iValue.Addr()
			}
			err := fv.Interface().(Relationship).Initialize(iValue.Interface(), ft)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
