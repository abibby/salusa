package slices

func Filter[T any](slice []T, cb func(T) bool) []T {
	result := []T{}
	for _, v := range slice {
		if cb(v) {
			result = append(result, v)
		}
	}
	return result
}
