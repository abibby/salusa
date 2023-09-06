package slices

type Equaler[T any] interface {
	Equal(T) bool
}

func Has[T comparable](haystack []T, needle T) bool {
	var iNeedle any = needle
	if eq, ok := iNeedle.(Equaler[T]); ok {
		for _, v := range haystack {
			if eq.Equal(v) {
				return true
			}
		}
	} else {
		for _, v := range haystack {
			if v == needle {
				return true
			}
		}
	}
	return false
}
