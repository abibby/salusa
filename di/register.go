package di

import (
	"context"
	"reflect"
)

func RegisterFactory(ctx context.Context, factory Factory) {
	dp := GetDependencyProvider(ctx)
	dp.Register(factory)
}

func Register[T any](ctx context.Context, factory func(ctx context.Context, tag string) (T, error)) {
	RegisterFactory(ctx, NewFactoryFunc(factory))
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
	RegisterFactory(ctx, NewSingletonFactory(factory()))
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
	RegisterFactory(ctx, NewLazySingletonFactory(func() (T, error) {
		return factory()
	}))
}
func RegisterValue(ctx context.Context, t reflect.Type, factory func(ctx context.Context, tag string) (reflect.Value, error)) {
	RegisterFactory(ctx, NewValueFactory(t, factory))
}

func (d *DependencyProvider) Register(factory Factory) {
	d.factories[factory.Type()] = factory
}
