package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/set"
)

type FillOptions struct {
	autoResolve set.Set[reflect.Type]
}

func newFillOptions() *FillOptions {
	return &FillOptions{
		autoResolve: set.Set[reflect.Type]{},
	}
}

type FillOption func(*FillOptions) *FillOptions

func (dp *DependencyProvider) Fill(ctx context.Context, v any, opts ...FillOption) error {
	opt := newFillOptions()
	for _, o := range opts {
		opt = o(opt)
	}
	return dp.fill(ctx, reflect.ValueOf(v), opt)
}
func (dp *DependencyProvider) fill(ctx context.Context, v reflect.Value, opt *FillOptions) error {
	if opt == nil {
		opt = newFillOptions()
	}
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("di: Fill(non-pointer "+v.Type().Name()+"): %w", ErrFillParameters)
	}

	if v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("di: Fill(non-struct "+v.Type().Name()+"): %w", ErrFillParameters)
	}

	return helpers.EachField(v, func(sf reflect.StructField, fv reflect.Value) error {
		tag, ok := sf.Tag.Lookup("inject")
		if !(ok || opt.autoResolve.Has(fv.Type())) {
			return nil
		}

		v, err := dp.resolve(ctx, sf.Type, tag)
		if err == nil {
			fv.Set(reflect.ValueOf(v))
			return nil
		} else if !errors.Is(err, ErrNotRegistered) {
			return fmt.Errorf("failed to fill: %w", err)
		}

		return fmt.Errorf("unable to fill field %s: %w", sf.Name, ErrNotRegistered)
	})
}

func AutoResolve[T any]() FillOption {
	return func(fo *FillOptions) *FillOptions {
		fo.autoResolve.Add(getType[T]())
		return fo
	}
}
