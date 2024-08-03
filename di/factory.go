package di

import (
	"context"
	"reflect"
	"sync"
)

type Factory interface {
	Build(ctx context.Context, tag string) (any, error)
	Type() reflect.Type
}
type Singleton interface {
	Factory
	Peek() (any, error, bool)
}

// ===============================
// ||                           ||
// ||        FactoryFunc        ||
// ||                           ||
// ===============================

type FactoryFunc[T any] func(ctx context.Context, tag string) (T, error)

func NewFactoryFunc[T any](f func(ctx context.Context, tag string) (T, error)) FactoryFunc[T] {
	return FactoryFunc[T](f)
}

var _ Factory = (FactoryFunc[any])(nil)

func (f FactoryFunc[T]) Build(ctx context.Context, tag string) (any, error) {
	return f(ctx, tag)
}
func (f FactoryFunc[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

// ================================
// ||                            ||
// ||        ValueFactory        ||
// ||                            ||
// ================================

type ValueFactory struct {
	factory func(ctx context.Context, tag string) (reflect.Value, error)
	typ     reflect.Type
}

func NewValueFactory(typ reflect.Type, factory func(ctx context.Context, tag string) (reflect.Value, error)) *ValueFactory {
	return &ValueFactory{
		factory: factory,
		typ:     typ,
	}
}

var _ Factory = (*ValueFactory)(nil)

func (f *ValueFactory) Build(ctx context.Context, tag string) (any, error) {
	v, err := f.factory(ctx, tag)
	return v.Interface(), err
}
func (f *ValueFactory) Type() reflect.Type {
	return f.typ
}

// ================================
// ||                            ||
// ||      SingletonFactory      ||
// ||                            ||
// ================================

type SingletonFactory[T any] struct {
	value T
}

var _ Singleton = (*SingletonFactory[any])(nil)

func NewSingletonFactory[T any](value T) *SingletonFactory[T] {
	return &SingletonFactory[T]{value: value}
}
func (f *SingletonFactory[T]) Build(ctx context.Context, tag string) (any, error) {
	return f.value, nil
}
func (f *SingletonFactory[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

func (f *SingletonFactory[T]) Peek() (any, error, bool) {
	return f.value, nil, true
}

// ================================
// ||                            ||
// ||    LazySingletonFactory    ||
// ||                            ||
// ================================

type LazySingletonFactory[T any] struct {
	once    sync.Once
	ready   bool
	factory func() (T, error)

	value any
	err   error
}

var _ Singleton = (*LazySingletonFactory[any])(nil)

func NewLazySingletonFactory[T any](factory func() (T, error)) *LazySingletonFactory[T] {
	return &LazySingletonFactory[T]{
		once:    sync.Once{},
		factory: factory,
	}
}

func (f *LazySingletonFactory[T]) Build(ctx context.Context, tag string) (any, error) {
	f.once.Do(f.Load)
	return f.value, f.err
}
func (f *LazySingletonFactory[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

func (f *LazySingletonFactory[T]) Load() {
	f.value, f.err = f.factory()
	f.ready = true
}
func (f *LazySingletonFactory[T]) Peek() (any, error, bool) {
	return f.value, f.err, f.ready
}
