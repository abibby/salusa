package dbtest

import (
	"reflect"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/relationship"
)

type Factory[T model.Model] func() T

func NewFactory[T model.Model](cb func() T) Factory[T] {
	return Factory[T](cb)
}

type CountFactory[T model.Model] struct {
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

func (f Factory[T]) Create(tx database.DB) T {
	m := f()
	err := model.Save(tx, m)
	if err != nil {
		panic(err)
	}
	err = relationship.InitializeRelationships(m)
	if err != nil {
		panic(err)
	}
	return m
}

func (f *CountFactory[T]) Create(tx database.DB) []T {
	models := make([]T, f.count)
	for i := 0; i < f.count; i++ {
		m := f.factory()
		err := model.Save(tx, m)
		if err != nil {
			panic(err)
		}
		models[i] = m
	}
	return models
}
