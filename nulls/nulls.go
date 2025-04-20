package nulls

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
)

var nullBytes = []byte("null")

type Nullable[T any] struct {
	value T
	valid bool
}

var _ json.Marshaler = Nullable[int]{}
var _ json.Unmarshaler = (*Nullable[int])(nil)
var _ sql.Scanner = (*Nullable[int])(nil)

func New[T any](n T) Nullable[T] {
	return Nullable[T]{
		value: n,
		valid: true,
	}
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{
		valid: false,
	}
}

func (n Nullable[T]) Val() T {
	return n.value
}

func (n Nullable[T]) String() string {
	if !n.valid {
		return ""
	}
	return fmt.Sprint(n.value)
}

func (n Nullable[T]) IsNull() bool {
	return !n.valid
}

func (n Nullable[T]) Ok() (T, bool) {
	return n.value, n.valid
}

// MarshalJSON implements json.Marshaler.
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return nullBytes, nil
	}
	return json.Marshal(n.value)
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		var zero T
		n.value = zero
		n.valid = false
		return nil
	}

	n.valid = true
	return json.Unmarshal(b, &n.value)
}

func (n *Nullable[T]) Scan(value any) error {
	v := reflect.ValueOf(value)
	if isNil(v) {
		var zero T
		n.value = zero
		n.valid = false
		return nil
	}

	t := reflect.TypeFor[T]()
	if !v.CanConvert(t) || (t.Kind() == reflect.String && v.Kind() != reflect.String) {
		return fmt.Errorf("Nullable.Scan: cannot scan type %T into Nullable[%T]", value, *new(T))
	}
	n.value = v.Convert(t).Interface().(T)
	n.valid = true
	return nil
}

func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.valid {
		return nil, nil
	}

	return n.value, nil
}

func isNil(v reflect.Value) bool {
	if (v == reflect.Value{}) {
		return true
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}
	return false
}
