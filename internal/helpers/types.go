package helpers

import "reflect"

func NewOf[T any]() T {
	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem()).Interface().(T)
	}
	return zero
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
