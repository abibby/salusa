package stream

import (
	"iter"
	"slices"
)

type Stream[T any] struct {
	seq iter.Seq[T]
}

func Of[T any](seq iter.Seq[T]) *Stream[T] {
	return &Stream[T]{
		seq: seq,
	}
}
func OfSlice[T any](s []T) *Stream[T] {
	return Of(func(yield func(T) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	})
}

func (s *Stream[T]) All() iter.Seq[T] {
	return s.seq
}

func (s *Stream[T]) Slice() []T {
	return slices.Collect(s.seq)
}

func (s *Stream[T]) Filter(fn func(T) bool) *Stream[T] {
	return Of(func(yield func(T) bool) {
		for v := range s.seq {
			if !fn(v) {
				continue
			}
			if !yield(v) {
				return
			}
		}
	})
}

func (s *Stream[T]) Limit(limit int) *Stream[T] {
	return Of(func(yield func(T) bool) {
		i := 0
		for v := range s.seq {
			if i >= limit {
				return
			}
			i++
			if !yield(v) {
				return
			}
		}
	})
}

func (s *Stream[T]) Skip(skip int) *Stream[T] {
	return Of(func(yield func(T) bool) {
		i := 0
		for v := range s.seq {
			i++
			if i <= skip {
				continue
			}
			if !yield(v) {
				return
			}
		}
	})
}

func (s *Stream[T]) Find(fn func(T) bool) (T, bool) {
	for v := range s.seq {
		if fn(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

func (s *Stream[T]) Sort(cmp func(a, b T) int) *Stream[T] {
	slice := s.Slice()
	slices.SortFunc(slice, cmp)
	return OfSlice(slice)
}
