package nulls

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

var nullBytes = []byte("null")

type Null[T any] sql.Null[T]

var _ json.Marshaler = Null[int]{}
var _ json.Unmarshaler = (*Null[int])(nil)
var _ sql.Scanner = (*Null[int])(nil)
var _ driver.Valuer = Null[int]{}

func New[T any](n T) Null[T] {
	return Null[T]{
		V:     n,
		Valid: true,
	}
}

// MarshalJSON implements json.Marshaler.
func (n Null[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return nullBytes, nil
	}
	return json.Marshal(n.V)
}

// UnmarshalJSON implements json.Unmarshaler.
func (n *Null[T]) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullBytes) {
		var zero T
		n.V = zero
		n.Valid = false
		return nil
	}

	n.Valid = true
	return json.Unmarshal(b, &n.V)
}

func (n *Null[T]) Scan(value any) error {
	return (*sql.Null[T])(n).Scan(value)
}

func (n Null[T]) Value() (driver.Value, error) {
	return sql.Null[T](n).Value()
}
