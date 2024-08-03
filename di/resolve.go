package di

import (
	"context"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
)

var (
	contextType            = reflect.TypeFor[context.Context]()
	dependencyProviderType = reflect.TypeFor[*DependencyProvider]()
)

func Resolve[T any](ctx context.Context) (T, error) {
	dp := GetDependencyProvider(ctx)

	var result T
	v, err := dp.resolve(ctx, reflect.TypeFor[T](), "", false, nil)
	if v != nil {
		result = v.(T)
	}
	return result, err
}

func ResolveFill[T any](ctx context.Context) (T, error) {
	dp := GetDependencyProvider(ctx)

	var result T
	v, err := dp.resolve(ctx, reflect.TypeFor[T](), "", true, nil)
	if v != nil {
		result = v.(T)
	}
	return result, err
}

func (dp *DependencyProvider) resolve(ctx context.Context, t reflect.Type, tag string, fill bool, opt *FillOptions) (any, error) {
	switch t {
	case contextType:
		return ctx, nil
	case dependencyProviderType:
		return dp, nil
	}

	f, ok := dp.factories[t]
	if ok {
		return f.Build(ctx, tag)
	}

	if !fill && !isFillable(t, tag) {
		return nil, errNotRegistered(t)
	}

	v := helpers.Create(t)
	err := dp.fill(ctx, v, opt)
	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}
