package di

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/abibby/salusa/internal/helpers"
)

type Factory interface {
	Build(ctx context.Context, dp *DependencyProvider, tag string) (any, error)
	Type() reflect.Type
}
type Singleton interface {
	Peek() (any, error, bool)
}

type Dependant interface {
	DependsOn() []reflect.Type
}

func dependancies(t reflect.Type) []reflect.Type {
	if !isFillable(t) {
		return []reflect.Type{t}
	}
	deps := []reflect.Type{}
	for _, sf := range helpers.GetFields(t.Elem()) {
		if _, ok := sf.Tag.Lookup("inject"); ok {
			deps = append(deps, sf.Type)
		}
	}
	return deps
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

func (f FactoryFunc[T]) Build(ctx context.Context, dp *DependencyProvider, tag string) (any, error) {
	return f(ctx, tag)
}
func (f FactoryFunc[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}

// ===============================
// ||                           ||
// ||      FactoryFuncWith      ||
// ||                           ||
// ===============================

type FactoryFuncWith[T, W any] func(ctx context.Context, tag string, with W) (T, error)

func NewFactoryFuncWith[T, W any](f func(ctx context.Context, tag string, with W) (T, error)) FactoryFuncWith[T, W] {
	return FactoryFuncWith[T, W](f)
}

var _ Factory = (FactoryFuncWith[any, any])(nil)

func (f FactoryFuncWith[T, W]) Build(ctx context.Context, dp *DependencyProvider, tag string) (any, error) {
	var with W
	err := dp.Fill(ctx, &with)
	if err != nil {
		var zero T
		return zero, err
	}
	return f(ctx, tag, with)
}
func (f FactoryFuncWith[T, W]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}
func (f FactoryFuncWith[T, W]) DependsOn() []reflect.Type {
	return dependancies(reflect.TypeFor[W]())
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

func (f *ValueFactory) Build(ctx context.Context, dp *DependencyProvider, tag string) (any, error) {
	v, err := f.factory(ctx, tag)
	if err == nil && v.Type() != f.typ {
		return reflect.Zero(f.typ), fmt.Errorf("invalid type %v expected %v", v.Type(), f.typ)
	}
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
func (f *SingletonFactory[T]) Build(ctx context.Context, dp *DependencyProvider, tag string) (any, error) {
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

	value T
	err   error
}

var _ Singleton = (*LazySingletonFactory[any])(nil)

func NewLazySingletonFactory[T any](factory func() (T, error)) *LazySingletonFactory[T] {
	return &LazySingletonFactory[T]{
		once:    sync.Once{},
		factory: factory,
	}
}

func (f *LazySingletonFactory[T]) Build(ctx context.Context, dp *DependencyProvider, tag string) (any, error) {
	f.once.Do(f.load)
	return f.value, f.err
}
func (f *LazySingletonFactory[T]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}
func (f *LazySingletonFactory[T]) load() {
	f.value, f.err = f.factory()
	f.ready = true
}
func (f *LazySingletonFactory[T]) Peek() (any, error, bool) {
	return f.value, f.err, f.ready
}

// ================================
// ||                            ||
// ||  LazySingletonWithFactory  ||
// ||                            ||
// ================================

type LazySingletonWithFactory[T, W any] struct {
	done    atomic.Uint32
	m       sync.Mutex
	factory func(with W) (T, error)

	value T
	err   error
}

var _ Singleton = (*LazySingletonWithFactory[any, any])(nil)

func NewLazySingletonWithFactory[T, W any](factory func(with W) (T, error)) *LazySingletonWithFactory[T, W] {
	return &LazySingletonWithFactory[T, W]{
		factory: factory,
	}
}

func (f *LazySingletonWithFactory[T, W]) Build(ctx context.Context, dp *DependencyProvider, tag string) (any, error) {
	// implanted matching sync.Once to allow for inlining the fast path
	if f.done.Load() == 0 {
		f.load(ctx, dp)
	}
	return f.value, f.err
}
func (f *LazySingletonWithFactory[T, W]) Type() reflect.Type {
	return reflect.TypeFor[T]()
}
func (f *LazySingletonWithFactory[T, W]) load(ctx context.Context, dp *DependencyProvider) {
	f.m.Lock()
	defer f.m.Unlock()
	if f.done.Load() == 0 {
		defer f.done.Store(1)

		var with W
		err := dp.Fill(ctx, &with)
		if err != nil {
			f.err = err
			return
		}
		f.value, f.err = f.factory(with)
	}
}
func (f *LazySingletonWithFactory[T, W]) Peek() (any, error, bool) {
	return f.value, f.err, f.done.Load() != 0
}
func (f *LazySingletonWithFactory[T, W]) DependsOn() []reflect.Type {
	return dependancies(reflect.TypeFor[W]())
}
