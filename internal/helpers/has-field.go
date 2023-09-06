package helpers

import "reflect"

func HasField(v any, field string) bool {
	return hasField(reflect.TypeOf(v), field)
}

func hasField(rt reflect.Type, key string) bool {
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	if rt.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if ft.Anonymous {
			ok := hasField(rt.Field(i).Type, key)
			if ok {
				return true
			}
			continue
		}
		if FieldName(ft) == key {
			return true
		}
	}
	return false
}
