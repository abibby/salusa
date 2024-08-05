package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/abibby/salusa/internal/helpers"
)

// var ErrNotFillable = errors.New("struct is not fillable")

func Fill(ctx context.Context, v any) error {
	dp := GetDependencyProvider(ctx)
	return dp.Fill(ctx, v)
}
func (dp *DependencyProvider) Fill(ctx context.Context, v any) error {
	return dp.fill(ctx, reflect.ValueOf(v), "")
}
func (dp *DependencyProvider) fill(ctx context.Context, v reflect.Value, tag string) error {
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("di: Fill(non-pointer "+v.Type().String()+"): %w", ErrFillParameters)
	}

	if v.IsNil() {
		return fmt.Errorf("di: Fill(nil)")
	}

	if ok, err := dp.resolve(ctx, v, tag); ok {
		return err
	}

	if !isFillable(v.Type()) {
		if isFillable(v.Elem().Type()) {
			v = v.Elem()
			v.Set(helpers.Create(v.Type()))
		} else {
			return errNotRegistered(v.Type())
		}
	}

	err := helpers.EachField(v, func(sf reflect.StructField, fv reflect.Value) error {
		if !sf.IsExported() {
			return nil
		}

		rawTag, ok := sf.Tag.Lookup("inject")
		if !ok {
			return nil
		}
		tag := parseTag(rawTag)

		nv := reflect.New(sf.Type)
		err := dp.fill(ctx, nv, tag.Name)
		if errors.Is(err, ErrNotRegistered) {
			if tag.Optional {
				return nil
			} else {
				return fmt.Errorf("unable to fill field %s.%s: %w", v.Type().String(), sf.Name, err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to fill: %w", err)
		}
		fv.Set(nv.Elem())
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func isFillable(t reflect.Type) bool {
	if t.Kind() != reflect.Pointer {
		return false
	}
	if t.Elem().Kind() != reflect.Struct {
		return false
	}
	for _, sf := range helpers.GetFields(t.Elem()) {
		if _, ok := sf.Tag.Lookup("inject"); ok {
			return true
		}
	}
	return false
}

func (dp *DependencyProvider) resolve(ctx context.Context, v reflect.Value, tag string) (bool, error) {
	f, ok := dp.factories[v.Type().Elem()]
	if !ok {
		return false, nil
	}

	res, err := f.Build(ctx, dp, tag)
	if err != nil {
		return true, err
	}
	v.Elem().Set(reflect.ValueOf(res))
	return true, nil
}

type fillTag struct {
	Name     string
	Optional bool
}

func parseTag(rawTag string) *fillTag {
	parts := strings.Split(rawTag, ",")
	tag := &fillTag{}
	tag.Name = parts[0]
	for _, p := range parts[1:] {
		if p == "optional" {
			tag.Optional = true
		}
	}
	return tag
}
