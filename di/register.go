package di

import (
	"context"
	"reflect"
)

func Register[T any](dp *DependencyProvider, factory DependencyFactory[T]) {
	dp.Register(factory)
}

func RegisterSingleton[T any](dp *DependencyProvider, factory func() T) {
	v := factory()
	dp.Register(func(ctx context.Context, tag string) (T, error) {
		return v, nil
	})
}

func RegisterLazySingleton[T any](dp *DependencyProvider, factory func() T) {
	var v T
	initialized := false
	dp.Register(func(ctx context.Context, tag string) (T, error) {
		if !initialized {
			initialized = true
			v = factory()
		}
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
		iErr := out[1].Interface()
		var err error
		if iErr != nil {
			err = iErr.(error)
		}
		return out[0].Interface(), err
	}
}
