package helpers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

var (
	ErrNotFound = errors.New("not found")
)

func GetValue(v any, key string) (any, bool) {
	result, err := RGetValue(reflect.ValueOf(v), key)
	if err != nil {
		spew.Dump(err)
	}
	return result, err == nil
}
func RGetValue(v reflect.Value, key string) (any, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, fmt.Errorf("v must not be nil")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("v must be a struct")
	}
	rt := v.Type()
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if ft.Anonymous {
			result, err := RGetValue(v.Field(i), key)
			if errors.Is(err, ErrNotFound) {
				continue
			} else if err != nil {
				return result, err
			}
			return result, nil
		}
		if FieldName(ft) == key {
			return v.Field(i).Interface(), nil
		}
	}
	return nil, ErrNotFound
}

func FieldName(f reflect.StructField) string {
	return DBTag(f).Name
}
