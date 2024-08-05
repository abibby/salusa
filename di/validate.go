package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/validate"
)

type DIValidator struct {
	dp  *DependencyProvider
	typ reflect.Type
}

var _ validate.Validator = (*DIValidator)(nil)

func (v *DIValidator) Validate(ctx context.Context) error {
	errs := []error{}
	t := v.typ
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	for _, sf := range helpers.GetFields(t) {
		if !sf.IsExported() {
			continue
		}

		_, ok := sf.Tag.Lookup("inject")
		if !ok {
			continue
		}

		switch sf.Type {
		case contextType, dependencyProviderType:
			continue
		}

		_, ok = v.dp.factories[sf.Type]
		if !ok {
			errs = append(errs, fmt.Errorf("missing dependancy %s on %s.%s", sf.Type, v.typ, sf.Name))
		}
	}
	return errors.Join(errs...)
}

func Validator(ctx context.Context, rootType reflect.Type) *DIValidator {
	return GetDependencyProvider(ctx).Validator(rootType)
}
func (dp *DependencyProvider) Validator(rootType reflect.Type) *DIValidator {
	return &DIValidator{
		dp:  dp,
		typ: rootType,
	}
}
