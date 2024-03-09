package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
)

type DependencyProvider struct {
	factories map[reflect.Type]func(ctx context.Context, tag string) (any, error)
}

type DependencyFactory[T any] func(ctx context.Context, tag string) (T, error)

var (
	contextType = getType[context.Context]()
	stringType  = getType[string]()
	errorType   = getType[error]()
)

var (
	ErrInvalidDependencyFactory = errors.New("dependency factories must match the type di.DependencyFactory")
	ErrNotRegistered            = errors.New("dependency not registered")
	ErrFillParameters           = errors.New("invalid fill parameters")
)

var defaultProvider = NewDependencyProvider()

func NewDependencyProvider() *DependencyProvider {
	return &DependencyProvider{
		factories: map[reflect.Type]func(ctx context.Context, tag string) (any, error){},
	}
}

func Register[T any](factory DependencyFactory[T]) {
	defaultProvider.Register(factory)
}

func RegisterSingleton[T any](factory func() T) {
	v := factory()
	defaultProvider.Register(func(ctx context.Context, tag string) (T, error) {
		return v, nil
	})
}

func RegisterLazySingleton[T any](factory func() T) {
	var v T
	initialized := false
	defaultProvider.Register(func(ctx context.Context, tag string) (T, error) {
		if !initialized {
			initialized = true
			v = factory()
		}
		return v, nil
	})
}

func (d *DependencyProvider) Register(factory any) {
	v := reflect.ValueOf(factory)
	t := v.Type()
	if t.Kind() != reflect.Func ||
		t.NumIn() != 2 ||
		t.In(0) != contextType ||
		t.In(1) != stringType ||
		t.NumOut() != 2 ||
		t.Out(1) != errorType {
		panic(ErrInvalidDependencyFactory)
	}

	d.factories[t.Out(0)] = func(ctx context.Context, tag string) (any, error) {
		out := v.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(tag),
		})
		iErr := out[1].Interface()
		var err error
		if iErr != nil {
			err = iErr.(error)
		}
		return out[0].Interface(), err
	}
}

func Fill(ctx context.Context, v any) error {
	return defaultProvider.Fill(ctx, v)
}
func (dp *DependencyProvider) Fill(ctx context.Context, v any) error {
	return dp.fill(ctx, reflect.ValueOf(v))
}
func (dp *DependencyProvider) fill(ctx context.Context, v reflect.Value) error {
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("di: Fill(non-pointer "+v.Type().Name()+"): %w", ErrFillParameters)
	}
	if v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("di: Fill(non-struct "+v.Type().Name()+"): %w", ErrFillParameters)
	}

	return helpers.EachField(v, func(sf reflect.StructField, fv reflect.Value) error {
		tag, ok := sf.Tag.Lookup("inject")
		if !ok {
			return nil
		}

		v, err := dp.resolve(sf.Type, ctx, tag)
		if err == nil {
			fv.Set(reflect.ValueOf(v))
			return nil
		} else if !errors.Is(err, ErrNotRegistered) {
			return fmt.Errorf("failed to fill: %w", err)
		}

		return fmt.Errorf("unable to fill field %s: %w", sf.Name, ErrNotRegistered)
	})
}

func Resolve[T any](ctx context.Context) (T, error) {
	return ResolveProvider[T](defaultProvider, ctx)
}

func ResolveProvider[T any](dp *DependencyProvider, ctx context.Context) (T, error) {
	var result T
	v, err := dp.resolve(getType[T](), ctx, "")
	if v != nil {
		result = v.(T)
	}
	return result, err
}
func (dp *DependencyProvider) resolve(t reflect.Type, ctx context.Context, tag string) (any, error) {
	if t == contextType {
		return ctx, nil
	}

	f, ok := dp.factories[t]
	if ok {
		return f(ctx, tag)
	}

	if t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Struct {
		return nil, ErrNotRegistered
	}
	v := reflect.New(t.Elem())
	err := dp.fill(ctx, v)
	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}

func getType[T any]() reflect.Type {
	var v T
	return reflect.TypeOf(&v).Elem()
}
