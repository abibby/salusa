package di

import (
	"context"
	"reflect"
	"sync"
)

type DependencyFactory[T any] func(ctx context.Context, tag string) (T, error)

var (
	contextType            = reflect.TypeFor[context.Context]()
	stringType             = reflect.TypeFor[string]()
	errorType              = reflect.TypeFor[error]()
	dependencyProviderType = reflect.TypeFor[*DependencyProvider]()
)

func Register[T any](ctx context.Context, factory DependencyFactory[T]) {
	dp := GetDependencyProvider(ctx)
	dp.Register(factory)
}
func RegisterWith[T, W any](ctx context.Context, factory func(ctx context.Context, tag string, with W) (T, error)) {
	Register(ctx, func(ctx context.Context, tag string) (T, error) {
		with, err := ResolveFill[W](ctx)
		if err != nil {
			var zero T
			return zero, err
		}
		return factory(ctx, tag, with)
	})
}

func RegisterSingleton[T any](ctx context.Context, factory func() T) {
	dp := GetDependencyProvider(ctx)
	value := factory()
	dp.Register(func(ctx context.Context, tag string) (T, error) {
		return value, nil
	})
}

func RegisterLazySingletonWith[T, W any](ctx context.Context, factory func(with W) (T, error)) {
	RegisterLazySingleton(ctx, func() (T, error) {
		with, err := ResolveFill[W](ctx)
		if err != nil {
			var zero T
			return zero, err
		}
		return factory(with)
	})
}

func RegisterLazySingleton[T any](ctx context.Context, factory func() (T, error)) {
	dp := GetDependencyProvider(ctx)
	onceFactory := sync.OnceValues(factory)
	dp.Register(func(ctx context.Context, tag string) (T, error) {
		return onceFactory()
	})
}
func RegisterValue(ctx context.Context, t reflect.Type, factory func(ctx context.Context, tag string) (reflect.Value, error)) {
	dp := GetDependencyProvider(ctx)
	dp.RegisterValue(t, factory)
}

// factory should be of type DependencyFactory[T any] func(ctx context.Context, tag string) (T, error)
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

	d.RegisterValue(t.Out(0), func(ctx context.Context, tag string) (reflect.Value, error) {
		out := v.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(tag),
		})
		err, _ := out[1].Interface().(error)
		return out[0], err
	})
}

func (d *DependencyProvider) RegisterValue(t reflect.Type, factory func(ctx context.Context, tag string) (reflect.Value, error)) {
	d.factories[t] = func(ctx context.Context, tag string) (any, error) {
		v, err := factory(ctx, tag)
		return v.Interface(), err
	}
}
