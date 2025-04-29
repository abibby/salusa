package sets

import (
	"iter"
)

type OrderedSet[T comparable] struct {
	slice []T
}

func NewOrderedSet[T comparable](values ...T) *OrderedSet[T] {
	s := &OrderedSet[T]{slice: make([]T, 0, len(values))}
	s.Add(values...)
	return s
}

func (s *OrderedSet[T]) Add(values ...T) {
	for _, v := range values {
		s.insert(v)
	}
}

func (s *OrderedSet[T]) Delete(values ...T) {
	for _, v := range values {
		s.delete(v)
	}
}
func (s *OrderedSet[T]) Has(value T) bool {
	return s.index(value) != -1
}
func (s *OrderedSet[T]) Clone() Set[T] {
	newSet := &OrderedSet[T]{
		slice: make([]T, len(s.slice)),
	}
	copy(newSet.slice, s.slice)
	return newSet
}
func (s *OrderedSet[T]) Len() int {
	return len(s.slice)
}
func (s *OrderedSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.slice {
			if !yield(v) {
				return
			}
		}
	}
}
func (s *OrderedSet[T]) Get(i int) T {
	return s.slice[i]
}

func (s *OrderedSet[T]) insert(value T) {
	if s.Has(value) {
		return
	}
	s.slice = append(s.slice, value)
}

func (s *OrderedSet[T]) delete(value T) {
	index := s.index(value)

	if index == -1 {
		return
	}

	copy(s.slice[index:], s.slice[index+1:])
	s.slice = s.slice[:len(s.slice)-1]
}

func (s *OrderedSet[T]) index(value T) int {
	for i, v := range s.slice {
		if v == value {
			return i
		}
	}
	return -1
}
