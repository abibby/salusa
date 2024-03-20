package di

import (
	"context"
	"reflect"
)

func Resolve[T any](ctx context.Context, dp *DependencyProvider) (T, error) {
	var result T
	v, _, err := dp.resolve(ctx, getType[T](), "", false, nil)
	if v != nil {
		result = v.(T)
	}
	return result, err
}
func (dp *DependencyProvider) resolve(ctx context.Context, t reflect.Type, tag string, forceFill bool, opt *FillOptions) (any, Closer, error) {
	if t == contextType {
		return ctx, nil, nil
	}

	if t == dependencyProviderType {
		return dp, nil, nil
	}

	f, ok := dp.factories[t]
	if ok {
		return f(ctx, tag)
	}

	if !forceFill && !isFillable(t) {
		return nil, nil, ErrNotRegistered
	}
	v := reflectNew(t)
	closer, err := dp.fill(ctx, v, opt)
	if err != nil {
		return nil, nil, err
	}

	return v.Interface(), closer, nil
}
