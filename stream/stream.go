package stream

import (
	"iter"
	"slices"
)

type Stream[T any] struct {
	iterable iterable[T]
}

type iterable[T any] interface {
	All() iter.Seq[T]
}

type iterableSeq[T any] iter.Seq[T]

func (i iterableSeq[T]) All() iter.Seq[T] {
	return iter.Seq[T](i)
}

type iterableSlice[T any] struct {
	slice []T
}

func (i iterableSlice[T]) All() iter.Seq[T] {
	return slices.Values(i.slice)
}

func New[T any](seq iter.Seq[T]) *Stream[T] {
	return &Stream[T]{
		iterable: iterableSeq[T](seq),
	}
}
func Of[T any](s []T) *Stream[T] {
	return &Stream[T]{
		iterable: iterableSlice[T]{slice: s},
	}
}

func (s *Stream[T]) All() iter.Seq[T] {
	return s.iterable.All()
}

func (s *Stream[T]) Slice() []T {
	if i, ok := s.iterable.(iterableSlice[T]); ok {
		return i.slice
	}
	return slices.Collect(s.iterable.All())
}
