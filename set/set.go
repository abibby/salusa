package set

type Set[T comparable] map[T]struct{}

func New[T comparable]() Set[T] {
	return Set[T]{}
}
func (s Set[T]) Add(v T) {
	s[v] = struct{}{}
}
func (s Set[T]) Delete(v T) {
	delete(s, v)
}
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}
func (s Set[T]) Clone() Set[T] {
	newSet := Set[T]{}
	for v := range s {
		newSet.Add(v)
	}
	return newSet
}
func (s Set[T]) Len() int {
	return len(s)
}
