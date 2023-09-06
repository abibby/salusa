package dbtest

import (
	"reflect"

	"github.com/abibby/salusa/database/insert"
	"github.com/abibby/salusa/database/models"
	"github.com/abibby/salusa/internal/helpers"
)

type Factory[T models.Model] func() T

func NewFactory[T models.Model](cb func() T) Factory[T] {
	return Factory[T](cb)
}

type CountFactory[T models.Model] struct {
	factory Factory[T]
	count   int
}

func name[T any]() string {
	var m T
	return reflect.TypeOf(m).String()
}

func (f Factory[T]) Count(count int) *CountFactory[T] {
	return &CountFactory[T]{
		factory: f,
		count:   count,
	}
}
func (f Factory[T]) State(s func(T) T) Factory[T] {
	return func() T {
		return s(f())
	}
}

func (f Factory[T]) Create(tx helpers.QueryExecer) T {
	m := f()
	err := insert.Save(tx, m)
	if err != nil {
		panic(err)
	}
	return m
}

func (f *CountFactory[T]) Create(tx helpers.QueryExecer) []T {
	models := make([]T, f.count)
	for i := 0; i < f.count; i++ {
		m := f.factory()
		err := insert.Save(tx, m)
		if err != nil {
			panic(err)
		}
		models[i] = m
	}
	return models
}
