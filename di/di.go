package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type contextKey uint8

const (
	dpKey uint8 = iota
)

type DependencyProvider struct {
	factories map[reflect.Type]func(ctx context.Context, tag string) (any, error)
}

var (
	ErrInvalidDependencyFactory       = errors.New("dependency factories must match the type di.DependencyFactory")
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
		factories: map[reflect.Type]func(ctx context.Context, tag string) (any, error){},
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
