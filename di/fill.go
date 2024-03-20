package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/set"
	"github.com/davecgh/go-spew/spew"
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
	isFillablerType = getType[IsFillabler]()
)

func (dp *DependencyProvider) Fill(ctx context.Context, v any, opts ...FillOption) error {
	opt := newFillOptions()
	for _, o := range opts {
		opt = o(opt)
	}
	_, err := dp.fill(ctx, reflect.ValueOf(v), opt)
	return err
}
func (dp *DependencyProvider) fill(ctx context.Context, v reflect.Value, opt *FillOptions) (Closer, error) {
	if opt == nil {
		opt = newFillOptions()
	}
	if v.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("di: Fill(non-pointer "+v.Type().String()+"): %w", ErrFillParameters)
	}

	if v.IsNil() {
		return nil, fmt.Errorf("di: Fill(nil)")
	}

	if v.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("di: Fill(non-struct "+v.Type().String()+"): %w", ErrFillParameters)
	}

	closers := []Closer{}

	err := helpers.EachField(v, func(sf reflect.StructField, fv reflect.Value) error {
		if !sf.IsExported() {
			return nil
		}

		tag, ok := sf.Tag.Lookup("inject")
		if !(ok || opt.autoResolve.Has(fv.Type())) {
			return nil
		}

		v, closer, err := dp.resolve(ctx, sf.Type, tag, false, opt)
		if err == nil {
			closers = append(closers, closer)
			fv.Set(reflect.ValueOf(v))
			return nil
		} else if !errors.Is(err, ErrNotRegistered) {
			return fmt.Errorf("failed to fill: %w", err)
		}

		return fmt.Errorf("unable to fill field %s: %w", sf.Name, ErrNotRegistered)
	})
	if err != nil {
		return nil, err
	}

	spew.Dump(closers)

	return func(err error) error {
		closeErrors := []error{}
		for _, c := range closers {
			closeErr := c(err)
			if closeErr != nil {
				closeErrors = append(closeErrors, closeErr)
			}
		}
		if len(closeErrors) > 0 {
			return errors.Join(closeErrors...)
		}
		return nil
	}, nil
}

func AutoResolve[T any]() FillOption {
	return func(fo *FillOptions) *FillOptions {
		fo.autoResolve.Add(getType[T]())
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
