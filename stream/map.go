package stream

func (s *Stream[T]) Map[R any](fn func(T) R) *Stream[R] {
	return Of(func(yield func(R) bool) {
		for v := range s.seq {
			if !yield(fn(v)) {
				return
			}
		}
	})
}

func (s *Stream[T]) FlatMap[R any](fn func(T) []R) *Stream[R] {
	return Of(func(yield func(R) bool) {
		for a := range s.seq {
			for _, v := range fn(a) {
				if !yield(v) {
					return
				}
			}
		}
	})
}

func (s *Stream[T]) Reduce[R any](fn func(R, T) R) R {
	var accumulator R
	for a := range s.seq {
		accumulator = fn(accumulator, a)
	}
	return accumulator
}
