package slices

func Find[T comparable](haystack []T, compare func(v T) bool) (T, bool) {
	for _, v := range haystack {
		if compare(v) {
			return v, true
		}
	}

	var zero T

	return zero, false
}
