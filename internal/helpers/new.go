package helpers

import "reflect"

func NewOf[T any]() T {
	var a T
	t := reflect.TypeOf(a)
	if t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem()).Interface().(T)
	}
	return reflect.New(t).Elem().Interface().(T)
}
