package helpers

import (
	"fmt"
	"reflect"
)

var (
	ErrExpectedStruct = fmt.Errorf("expected a struct")
)

func EachField(v reflect.Value, cb func(sf reflect.StructField, fv reflect.Value) error) error {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fv := v.Field(i)
		if sf.Anonymous {
			err := EachField(fv, cb)
			if err != nil {
				return err
			}
		} else {
			err := cb(sf, fv)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
