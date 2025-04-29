package sets

import (
	"iter"
)

type Set[T comparable] interface {
	Add(v ...T)
	Delete(v ...T)
	Has(v T) bool
	Len() int
	All() iter.Seq[T]
	Clone() Set[T]
}

func New[T comparable](values ...T) Set[T] {
	return NewMapSet[T](values...)
}
