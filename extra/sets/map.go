package sets

import "iter"

type MapSet[T comparable] map[T]struct{}

func NewMapSet[T comparable](values ...T) MapSet[T] {
	s := MapSet[T]{}
	s.Add(values...)
	return s
}
func (s MapSet[T]) Add(values ...T) {
	for _, v := range values {
		s[v] = struct{}{}
	}
}
func (s MapSet[T]) Delete(values ...T) {
	for _, v := range values {
		delete(s, v)
	}
}
func (s MapSet[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}
func (s MapSet[T]) Clone() Set[T] {
	newSet := MapSet[T]{}
	for v := range s {
		newSet.Add(v)
	}
	return newSet
}
func (s MapSet[T]) Len() int {
	return len(s)
}
func (s MapSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if !yield(v) {
				return
			}
		}
	}
}
