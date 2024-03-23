package di

import (
	"context"
	"reflect"
	"sync"

	"github.com/abibby/salusa/internal/helpers"
)

type DependencyFactory[T any] func(ctx context.Context, tag string) (T, error)

var (
	contextType            = helpers.GetType[context.Context]()
	stringType             = helpers.GetType[string]()
	errorType              = helpers.GetType[error]()
	dependencyProviderType = helpers.GetType[*DependencyProvider]()
)

func Register[T any](dp *DependencyProvider, factory DependencyFactory[T]) {
	dp.Register(factory)
}
func RegisterWith[T, W any](dp *DependencyProvider, factory func(ctx context.Context, tag string, with W) (T, error)) {
	Register(dp, func(ctx context.Context, tag string) (T, error) {
		with, err := dp.resolve(ctx, helpers.GetType[W](), tag, nil)
		if err != nil {
			var zero T
			return zero, err
		}
		return factory(ctx, tag, with.(W))
	})
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

	d.factories[t.Out(0)] = func(ctx context.Context, tag string) (any, error) {
		out := v.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(tag),
		})
		err, _ := out[1].Interface().(error)
		return out[0].Interface(), err
	}
}
