package di

import (
	"context"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
)

func Resolve[T any](ctx context.Context, dp *DependencyProvider) (T, error) {
	var result T
	v, err := dp.resolve(ctx, helpers.GetType[T](), "", nil)
	if v != nil {
		result = v.(T)
	}
	return result, err
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

	if !isFillable(t) {
		return nil, ErrNotRegistered
	}
	v := reflectNew(t)
	err := dp.fill(ctx, v, opt)
	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}
