package sets

import (
	"cmp"
	"iter"
	"slices"
)

type SliceSet[T cmp.Ordered] struct {
	slice []T
}

func NewSliceSet[T cmp.Ordered](values ...T) *SliceSet[T] {
	s := &SliceSet[T]{}
	s.Add(values...)
	return s
}
func (s *SliceSet[T]) Add(values ...T) {
	for _, v := range values {
		s.insert(v)
	}
}

func (s *SliceSet[T]) Delete(values ...T) {
	for _, v := range values {
		s.delete(v)
	}
}
func (s *SliceSet[T]) Has(v T) bool {
	_, found := slices.BinarySearch(s.slice, v)
	return found
}
func (s *SliceSet[T]) Clone() Set[T] {
	newSet := &SliceSet[T]{
		slice: make([]T, len(s.slice)),
	}
	copy(newSet.slice, s.slice)
	return newSet
}
func (s *SliceSet[T]) Len() int {
	return len(s.slice)
}
func (s *SliceSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.slice {
			if !yield(v) {
				return
			}
		}
	}
}
func (s *SliceSet[T]) Get(i int) T {
	return s.slice[i]
}

func (s *SliceSet[T]) insert(value T) {
	idx, found := slices.BinarySearch(s.slice, value)
	if found {
		return
	}

	var zero T
	s.slice = append(s.slice, zero)

	copy(s.slice[idx+1:], s.slice[idx:])

	s.slice[idx] = value
}

func (s *SliceSet[T]) delete(value T) {
	index, found := slices.BinarySearch(s.slice, value)

	if !found {
		return
	}

	copy(s.slice[index:], s.slice[index+1:])
	s.slice = s.slice[:len(s.slice)-1]
}
