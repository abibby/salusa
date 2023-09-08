package helpers

import "reflect"

func Each(v any, cb func(v reflect.Value, pointer bool) error) error {
	pointer := false
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
		pointer = true
	}
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			err := Each(rv.Index(i).Interface(), cb)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if rv.Kind() != reflect.Struct {
		return nil
	}

	return cb(rv, pointer)
}
