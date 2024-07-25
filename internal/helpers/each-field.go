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

func GetFields(t reflect.Type) []reflect.StructField {
	return appendFields(make([]reflect.StructField, 0, t.NumField()), t)
}
func appendFields(fields []reflect.StructField, t reflect.Type) []reflect.StructField {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fields
	}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}
		if sf.Anonymous {
			fields = appendFields(fields, sf.Type)
			continue
		}

		fields = append(fields, sf)
	}
	return fields
}
