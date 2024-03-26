package di

import (
	"context"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
)

func Resolve[T any](ctx context.Context) (T, error) {
	dp := GetDependencyProvider(ctx)

	var result T
	v, err := dp.resolve(ctx, helpers.GetType[T](), "", nil)
	if v != nil {
		result = v.(T)
	}
	return result, err
}

func ResolveFill[T any](ctx context.Context) (T, error) {
	v := helpers.NewOf[T]()
	err := Fill(ctx, v)
	return v, err
}
func (dp *DependencyProvider) resolve(ctx context.Context, t reflect.Type, tag string, opt *FillOptions) (any, error) {
	if t == contextType {
		return ctx, nil
	}

	if t == dependencyProviderType {
		return dp, nil
	}

	f, ok := dp.factories[t]
	if ok {
		return f(ctx, tag)
	}

	if !isFillable(t, tag) {
		return nil, errNotRegistered(t)
	}
	v := reflectNew(t)
	err := dp.fill(ctx, v, opt)
	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}
