package maps

import "iter"

type Map[K, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Remove(key K)
	All() iter.Seq2[K, V]
}
