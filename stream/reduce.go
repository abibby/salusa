package stream

func (s *Stream[T]) Reduce[R any](fn func(R, T) R) R {
	var accumulator R
	return s.ReduceFrom(accumulator, fn)
}

func (s *Stream[T]) ReduceFrom[R any](accumulator R, fn func(R, T) R) R {
	for a := range s.iterable.All() {
		accumulator = fn(accumulator, a)
	}
	return accumulator
}
