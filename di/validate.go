package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
)

func Validate(ctx context.Context, v any, opts ...FillOption) error {
	return GetDependencyProvider(ctx).Validate(ctx, v, opts...)
}
func (dp *DependencyProvider) Validate(ctx context.Context, v any, opts ...FillOption) error {
	opt := newFillOptions()
	for _, o := range opts {
		opt = o(opt)
	}
	errs := []error{}

	err := helpers.EachField(reflect.ValueOf(v), func(sf reflect.StructField, fv reflect.Value) error {
		if !sf.IsExported() {
			return nil
		}

		_, ok := sf.Tag.Lookup("inject")
		if !ok && !opt.autoResolve.Has(fv.Type()) {
			return nil
		}

		switch sf.Type {
		case contextType, dependencyProviderType:
			return nil
		}

		_, ok = dp.factories[sf.Type]
		if !ok {
			errs = append(errs, fmt.Errorf("missing dependancy %s", sf.Type))
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
	return errors.Join(errs...)
}
