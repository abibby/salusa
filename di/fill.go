package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/set"
)

type IsFillabler interface {
	IsFillable() bool
}

type Fillable struct{}

func (*Fillable) IsFillable() bool {
	return true
}

type FillOptions struct {
	autoResolve set.Set[reflect.Type]
}

func newFillOptions() *FillOptions {
	return &FillOptions{
		autoResolve: set.Set[reflect.Type]{},
	}
}

type FillOption func(*FillOptions) *FillOptions

var (
	isFillablerType = helpers.GetType[IsFillabler]()
)

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
		return fmt.Errorf("di: Fill(non-pointer "+v.Type().String()+"): %w", ErrFillParameters)
	}

	if v.IsNil() {
		return fmt.Errorf("di: Fill(nil)")
	}

	if v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("di: Fill(non-struct "+v.Type().String()+"): %w", ErrFillParameters)
	}

	err := helpers.EachField(v, func(sf reflect.StructField, fv reflect.Value) error {
		if !sf.IsExported() {
			return nil
		}

		tag, ok := sf.Tag.Lookup("inject")
		if !ok && !opt.autoResolve.Has(fv.Type()) {
			return nil
		}

		v, err := dp.resolve(ctx, sf.Type, tag, opt)
		if errors.Is(err, ErrNotRegistered) {
			for t, _ := range dp.factories {
				fmt.Printf("%s\n", t)
			}
			return fmt.Errorf("unable to fill field %s (%s): %w", sf.Name, sf.Type.String(), ErrNotRegistered)
		} else if err != nil {
			return fmt.Errorf("failed to fill: %w", err)
		}

		fv.Set(reflect.ValueOf(v))
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func AutoResolve[T any]() FillOption {
	return func(fo *FillOptions) *FillOptions {
		fo.autoResolve.Add(helpers.GetType[T]())
		return fo
	}
}

func reflectNew(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem())
	}
	return reflect.New(t).Elem()
}

func isFillable(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer &&
		t.Elem().Kind() == reflect.Struct &&
		t.Implements(isFillablerType)
}
