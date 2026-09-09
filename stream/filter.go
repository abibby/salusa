package stream

func (s *Stream[T]) Filter(fn func(T) bool) *Stream[T] {
	return New(func(yield func(T) bool) {
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

func (s *Stream[T]) Find(fn func(T) bool) (T, bool) {
	for v := range s.seq {
		if fn(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}
