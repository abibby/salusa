package stream

func (s *Stream[T]) Limit(limit int) *Stream[T] {
	return New(func(yield func(T) bool) {
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
	return New(func(yield func(T) bool) {
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
