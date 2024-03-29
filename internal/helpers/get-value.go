package helpers

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotFound = errors.New("not found")
)

func GetValue(v any, key string) (any, bool) {
	result, err := RGetValue(reflect.ValueOf(v), key)
	if err != nil {
		return nil, false
	}
	return result.Interface(), err == nil
}
func RGetValue(v reflect.Value, key string) (reflect.Value, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.Value{}, fmt.Errorf("v must not be nil")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("v must be a struct")
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
			return v.Field(i), nil
		}
	}
	return reflect.Value{}, ErrNotFound
}

func FieldName(f reflect.StructField) string {
	return DBTag(f).Name
}
