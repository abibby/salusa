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
)

var defaultProvider = NewDependencyProvider()

func SetDefaultProvider(dp *DependencyProvider) {
	defaultProvider = dp
}

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

func Fill(ctx context.Context, v any, options ...any) error {
	return defaultProvider.Fill(ctx, v, options...)
}
func (dp *DependencyProvider) Fill(ctx context.Context, v any, options ...any) error {
	return dp.fill(ctx, reflect.ValueOf(v), options...)
}
func (dp *DependencyProvider) fill(ctx context.Context, v reflect.Value, options ...any) error {
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

		if sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Struct {
			v := reflect.New(sf.Type.Elem())
			err := dp.fill(ctx, v, options...)
			if err != nil {
				return err
			}
			fv.Set(v)
			return nil
		}
		return fmt.Errorf("unable to fill field %s: %w", sf.Name, ErrNotRegistered)
	})
}

func Resolve[T any](ctx context.Context) (T, error) {
	return ResolveProvider[T](defaultProvider, ctx)
}

func ResolveProvider[T any](dp *DependencyProvider, ctx context.Context) (T, error) {
	v, ok := dp.resolve(getType[T](), ctx, "")
	if v == nil {
		var zero T
		return zero, ok
	}
	return v.(T), ok
}
func (dp *DependencyProvider) resolve(t reflect.Type, ctx context.Context, tag string) (any, error) {
	if t == contextType {
		return ctx, nil
	}

	f, ok := dp.factories[t]
	if !ok {
		return nil, ErrNotRegistered
	}
	return f(ctx, tag)
}

func getType[T any]() reflect.Type {
	var v T
	return reflect.TypeOf(&v).Elem()
}
