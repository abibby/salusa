package di

import (
	"context"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
)

type DependancyProvider struct {
	factories map[reflect.Type]func(ctx context.Context, tag string) any
}

type DependancyFactory[T any] func(ctx context.Context, tag string) T

var (
	contextType = getType[context.Context]()
	stringType  = getType[string]()
)

var defaultProvider = NewDependamcyProvider()

func SetDefaultProvider(dp *DependancyProvider) {
	defaultProvider = dp
}

func NewDependamcyProvider() *DependancyProvider {
	return &DependancyProvider{
		factories: map[reflect.Type]func(ctx context.Context, tag string) any{},
	}
}

func Register[T any](factory DependancyFactory[T]) {
	defaultProvider.Register(factory)
}

func RegisterSinglton[T any](factory func() T) {
	v := factory()
	defaultProvider.Register(func(ctx context.Context, tag string) T {
		return v
	})
}

func RegisterLazySinglton[T any](factory func() T) {
	var v T
	initialized := false
	defaultProvider.Register(func(ctx context.Context, tag string) T {
		if !initialized {
			v = factory()
		}
		return v
	})
}

func (d *DependancyProvider) Register(factory any) {
	v := reflect.ValueOf(factory)
	t := v.Type()
	if t.Kind() != reflect.Func ||
		t.NumIn() != 2 ||
		t.In(0) != contextType ||
		t.In(1) != stringType ||
		t.NumOut() != 1 {
		panic("dependancy factories must match the type di.DependancyFactory")
	}

	d.factories[t.Out(0)] = func(ctx context.Context, tag string) any {
		out := v.Call([]reflect.Value{
			reflect.ValueOf(ctx),
			reflect.ValueOf(tag),
		})
		return out[0].Interface()
	}
}

func Fill(ctx context.Context, v any, options ...any) error {
	return defaultProvider.Fill(ctx, v, options...)
}
func (dp *DependancyProvider) Fill(ctx context.Context, v any, options ...any) error {
	return dp.fill(ctx, reflect.ValueOf(v), options...)
}
func (dp *DependancyProvider) fill(ctx context.Context, v reflect.Value, options ...any) error {
	return helpers.EachField(v, func(sf reflect.StructField, fv reflect.Value) error {
		tag, ok := sf.Tag.Lookup("inject")
		if !ok {
			return nil
		}

		v, ok := dp.resolve(sf.Type, ctx, tag)
		if ok {
			fv.Set(reflect.ValueOf(v))
			return nil
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
		return nil
	})
}

func Resolve[T any](ctx context.Context) (T, bool) {
	return ResolveProvider[T](defaultProvider, ctx)
}

func ResolveProvider[T any](dp *DependancyProvider, ctx context.Context) (T, bool) {
	v, ok := dp.resolve(getType[T](), ctx, "")
	if v == nil {
		var zero T
		return zero, ok
	}
	return v.(T), ok
}
func (dp *DependancyProvider) resolve(t reflect.Type, ctx context.Context, tag string) (any, bool) {
	if t == contextType {
		return ctx, true
	}

	f, ok := dp.factories[t]
	if !ok {
		return nil, false
	}
	return f(ctx, tag), true
}

func getType[T any]() reflect.Type {
	var v T
	return reflect.TypeOf(&v).Elem()
}
