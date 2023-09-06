package helpers

import (
	"reflect"
)

func GetValue(v any, key string) (any, bool) {
	return RGetValue(reflect.ValueOf(v), key)
}
func RGetValue(rv reflect.Value, key string) (any, bool) {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, false
	}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if ft.Anonymous {
			result, ok := RGetValue(rv.Field(i), key)
			if ok {
				return result, true
			}
			continue
		}
		if FieldName(ft) == key {
			return rv.Field(i).Interface(), true
		}
	}
	return nil, false
}

func FieldName(f reflect.StructField) string {
	return DBTag(f).Name
}
