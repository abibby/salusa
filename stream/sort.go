package stream

import (
	"slices"
)

func (s *Stream[T]) Sort(cmp func(a, b T) int) *Stream[T] {
	slice := s.Slice()
	slices.SortFunc(slice, cmp)
	return Of(slice)
}
