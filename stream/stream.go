package stream

import (
	"iter"
	"slices"
)

type Stream[T any] struct {
	seq iter.Seq[T]
}

func New[T any](seq iter.Seq[T]) *Stream[T] {
	return &Stream[T]{
		seq: seq,
	}
}
func Of[T any](s []T) *Stream[T] {
	return New(slices.Values(s))
}

func (s *Stream[T]) All() iter.Seq[T] {
	return s.seq
}

func (s *Stream[T]) Slice() []T {
	return slices.Collect(s.seq)
}
