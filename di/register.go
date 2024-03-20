package di

import (
	"context"
	"reflect"
	"sync"
)

type DependencyFactory[T any] func(ctx context.Context, tag string) (T, error)
type CloserDependencyFactory[T any] func(ctx context.Context, tag string) (T, Closer, error)

var (
	contextType            = getType[context.Context]()
	stringType             = getType[string]()
	errorType              = getType[error]()
	closerType             = getType[Closer]()
	dependencyProviderType = getType[*DependencyProvider]()
)

func Register[T any](dp *DependencyProvider, factory DependencyFactory[T]) {
	dp.Register(factory)
}
func RegisterCloser[T any](dp *DependencyProvider, factory CloserDependencyFactory[T]) {
	dp.RegisterCloser(factory)
}

func RegisterSingleton[T any](dp *DependencyProvider, factory func() T) {
	v := factory()
	dp.Register(func(ctx context.Context, tag string) (T, error) {
		return v, nil
	})
}

func RegisterLazySingleton[T any](dp *DependencyProvider, factory func() T) {
	var v T
	initialize := sync.OnceFunc(func() {
		v = factory()
	})
	dp.Register(func(ctx context.Context, tag string) (T, error) {
		initialize()
		return v, nil
	})
}

func (d *DependencyProvider) Register(factory any) {
	v := reflect.ValueOf(factory)
	t := v.Type()
	if t.Kind() != reflect.Func ||
		t.NumIn() != 2 ||
		t.In(0) != contextType ||
		t.In(1) != stringType ||
		t.NumOut() != 2 ||
		t.Out(1) != errorType {
		panic(ErrInvalidDependencyFactory)
	}

	d.factories[t.Out(0)] = func(ctx context.Context, tag string) (any, Closer, error) {
		out := v.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(tag),
		})
		err, _ := out[1].Interface().(error)
		closer := func(err error) error { return nil }
		return out[0].Interface(), closer, err
	}
}

func (d *DependencyProvider) RegisterCloser(factory any) {
	v := reflect.ValueOf(factory)
	t := v.Type()
	if t.Kind() != reflect.Func ||
		t.NumIn() != 2 ||
		t.In(0) != contextType ||
		t.In(1) != stringType ||
		t.NumOut() != 3 ||
		t.Out(1) != closerType ||
		t.Out(2) != errorType {
		panic(ErrInvalidDependencyFactory)
	}

	d.factories[t.Out(0)] = func(ctx context.Context, tag string) (any, Closer, error) {
		out := v.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(tag),
		})
		closer, _ := out[1].Interface().(Closer)
		err, _ := out[2].Interface().(error)
		return out[0].Interface(), closer, err
	}
}
