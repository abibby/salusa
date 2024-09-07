package helpers

import (
	"fmt"
	"reflect"
)

func NewOf[T any]() (T, error) {
	t := reflect.TypeFor[T]()
	v, err := RNewOf(t)
	if v.IsZero() {
		var zero T
		return zero, err
	}
	return v.Interface().(T), err
}
func RNewOf(t reflect.Type) (reflect.Value, error) {
	if t.Kind() == reflect.Interface {
		return reflect.Zero(t), fmt.Errorf("cannot create a new interface %s", t)
	}
	v := reflect.New(t).Elem()
	if t.Kind() == reflect.Pointer {
		v.Set(reflect.New(t.Elem()))
	}
	return v, nil
}
func CreateFor[T any]() reflect.Value {
	return Create(reflect.TypeFor[T]())
}

func Create(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem())
	}
	return reflect.New(t).Elem()
}

func Zero[T any]() T {
	var v T
	return v
}
