package di

import (
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
)

type DependancyProvider struct {
	factories map[reflect.Type]func() any
}

var DefaultProvider = NewDependamcyProvider()

func NewDependamcyProvider() *DependancyProvider {
	return &DependancyProvider{
		factories: map[reflect.Type]func() any{},
	}
}

type Registerer[T any] struct {
	dp          *DependancyProvider
	reflectType reflect.Type
}

func Register[T any]() *Registerer[T] {
	return RegisterProvider[T](DefaultProvider)
}
func RegisterProvider[T any](dp *DependancyProvider) *Registerer[T] {
	return &Registerer[T]{
		dp:          dp,
		reflectType: getType[T](),
	}
}

func (r *Registerer[T]) Singlton(instance T) {
	r.Factory(func() T {
		return instance
	})
}
func (r *Registerer[T]) Factory(f func() T) {
	r.dp.factories[getType[T]()] = func() any {
		return f()
	}
}

func Fill(v any, options ...any) error {
	return FillProvider(DefaultProvider, v, options...)
}
func FillProvider(dp *DependancyProvider, v any, options ...any) error {
	return helpers.EachField(reflect.ValueOf(v), func(sf reflect.StructField, fv reflect.Value) error {
		v, ok := resolvePReflect(dp, sf.Type)
		if !ok {
			return nil
		}
		fv.Set(reflect.ValueOf(v))
		return nil
	})
}

func Resolve[T any]() (T, bool) {
	return ResolveProvider[T](DefaultProvider)
}
func ResolveProvider[T any](dp *DependancyProvider) (T, bool) {
	v, ok := resolvePReflect(dp, getType[T]())
	if v == nil {
		var zero T
		return zero, ok
	}
	return v.(T), ok
}
func resolvePReflect(dp *DependancyProvider, t reflect.Type) (any, bool) {
	f, ok := dp.factories[t]
	if !ok {
		return nil, false
	}
	return f(), true
}

func getType[T any]() reflect.Type {
	var v T
	return reflect.TypeOf(&v).Elem()
}
