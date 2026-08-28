package stream

// Map returns a new Stream containing the results of applying the function fn
// to each element in the original Stream.
//
// The mapping function is evaluated lazily as the returned Stream is iterated.
// If the target sequence iteration is halted early by the consumer, mapping
// stops immediately.
func (s *Stream[T]) Map[R any](fn func(T) R) *Stream[R] {
	return New(func(yield func(R) bool) {
		for v := range s.seq {
			if !yield(fn(v)) {
				return
			}
		}
	})
}

func (s *Stream[T]) FlatMap[R any](fn func(T) []R) *Stream[R] {
	return New(func(yield func(R) bool) {
		for a := range s.seq {
			for _, v := range fn(a) {
				if !yield(v) {
					return
				}
			}
		}
	})
}
