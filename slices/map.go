package slices

func Map[T, U any](slice []T, cb func(v T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = cb(v)
	}
	return result
}
