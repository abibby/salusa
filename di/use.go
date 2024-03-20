package di

import (
	"context"
	"reflect"
)

func Use[T any](ctx context.Context, dp *DependencyProvider, cb func(v T) error, opts ...FillOption) error {
	return dp.Use(ctx, cb, opts...)
}

func (dp *DependencyProvider) Use(ctx context.Context, cb any, opts ...FillOption) error {
	opt := newFillOptions()
	for _, o := range opts {
		opt = o(opt)
	}

	v := reflect.ValueOf(cb)
	t := v.Type()
	if t.Kind() != reflect.Func ||
		t.NumIn() != 1 ||
		t.NumOut() != 1 ||
		t.Out(0) != errorType {
		panic(ErrInvalidDependencyFactory)
	}
	resolved, close, err := dp.resolve(ctx, t.In(0), "", true, opt)
	if err != nil {
		return err
	}
	result := v.Call([]reflect.Value{
		reflect.ValueOf(resolved),
	})
	if !result[0].IsNil() {
		err = result[0].Interface().(error)
	}
	closeErr := close(err)
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}

	return nil
}
