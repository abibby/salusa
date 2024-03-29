package helpers

import (
	"fmt"
	"reflect"
)

func NewOf[T any]() (T, error) {
	var zero [0]T
	t := reflect.TypeOf(zero).Elem()
	if t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem()).Interface().(T), nil
	}

	if t.Kind() == reflect.Interface {
		var zero T
		return zero, fmt.Errorf("cannot create a new interface %s", t)
	}
	return reflect.New(t).Elem().Interface().(T), nil
}
func Create(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem())
	}
	return reflect.New(t).Elem()
}

func GetType[T any]() reflect.Type {
	var v [0]T
	return reflect.TypeOf(v).Elem()
}

func Zero[T any]() T {
	var v T
	return v
}
