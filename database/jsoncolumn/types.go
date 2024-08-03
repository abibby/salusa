package jsoncolumn

import (
	"database/sql/driver"
)

type Map[K comparable, V any] map[K]V

func (j *Map[K, V]) Scan(src any) error {
	return Scan(j, src)
}
func (j Map[K, V]) Value() (driver.Value, error) {
	return Value(j)
}

type Slice[V any] []V

func (j *Slice[V]) Scan(src any) error {
	return Scan(j, src)
}
func (j Slice[V]) Value() (driver.Value, error) {
	return Value(j)
}
