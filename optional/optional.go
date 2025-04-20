package optional

type Option[T any] struct {
	Value T
	Valid bool
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		Value: v,
		Valid: true,
	}
}
func None[T any]() Option[T] {
	return Option[T]{}
}

// UnmarshalJSON implements json.Unmarshaler.
func (o Option[T]) UnmarshalJSON([]byte) error {
	panic("unimplemented")
}

// MarshalJSON implements json.Marshaler.
func (o Option[T]) MarshalJSON() ([]byte, error) {
	panic("unimplemented")
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (o Option[T]) UnmarshalBinary(data []byte) error {
	panic("unimplemented")
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (o Option[T]) MarshalBinary() (data []byte, err error) {
	panic("unimplemented")
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (o Option[T]) UnmarshalText(text []byte) error {
	panic("unimplemented")
}

// MarshalText implements encoding.TextMarshaler.
func (o Option[T]) MarshalText() (text []byte, err error) {
	panic("unimplemented")
}
