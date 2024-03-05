package di

type DependancyProvider struct {
	factories map[string]func() any
}

var defaultProvider = NewDependamcyProvider()

func NewDependamcyProvider() *DependancyProvider {
	return &DependancyProvider{
		factories: map[string]func() any{},
	}
}

type Registerer[T any] struct {
	dp     *DependancyProvider
	typeID string
}

func Register[T any]() *Registerer[T] {
	return RegisterP[T](defaultProvider)
}
func RegisterP[T any](dp *DependancyProvider) *Registerer[T] {
	return &Registerer[T]{
		dp:     dp,
		typeID: "something from the reflected type",
	}
}

func (r *Registerer[T]) Singlton(instance T) {
	r.Factory(func() T {
		return instance
	})
}
func (r *Registerer[T]) Factory(f func() T) {
	r.dp.factories[r.typeID] = func() any {
		return f()
	}
}

func Fill(v any) error {
	return FillP(defaultProvider, v)
}
func FillP(dp *DependancyProvider, v any) error {
	return nil
}

func Resolve[T any]() T {
	return ResolveP[T](defaultProvider)
}
func ResolveP[T any](dp *DependancyProvider) T {
	var zero T
	return zero
}
