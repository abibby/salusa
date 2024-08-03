package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type contextKey uint8

const (
	dpKey contextKey = iota
)

type DependencyProvider struct {
	factories map[reflect.Type]Factory
}

var (
	ErrNotRegistered                  = errors.New("dependency not registered")
	ErrFillParameters                 = errors.New("invalid fill parameters")
	ErrDependencyProviderNotInContext = errors.New("DependencyProvider not in context")
)

func errNotRegistered(t reflect.Type) error {
	return fmt.Errorf("%w: %s", ErrNotRegistered, t)
}

var defaultProvider = NewDependencyProvider()

func NewDependencyProvider() *DependencyProvider {
	return &DependencyProvider{
		factories: map[reflect.Type]Factory{},
	}
}

func ContextWithDependencyProvider(ctx context.Context, dp *DependencyProvider) context.Context {
	return context.WithValue(ctx, dpKey, dp)
}

func GetDependencyProvider(ctx context.Context) *DependencyProvider {
	v := ctx.Value(dpKey)
	if v == nil {
		return defaultProvider
	}
	dp, ok := v.(*DependencyProvider)
	if !ok {
		return defaultProvider
	}
	return dp
}
func TestDependencyProviderContext() context.Context {
	return ContextWithDependencyProvider(
		context.Background(),
		NewDependencyProvider(),
	)
}

func (dp *DependencyProvider) Singletons() []Singleton {
	singletons := make([]Singleton, 0, len(dp.factories))
	for _, f := range dp.factories {
		if s, ok := f.(Singleton); ok {
			singletons = append(singletons, s)
		}
	}
	return singletons
}
func Singletons(ctx context.Context) []Singleton {
	dp := GetDependencyProvider(ctx)
	return dp.Singletons()
}
