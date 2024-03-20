package di

import (
	"context"
	"errors"
	"reflect"
)

type Closer func(error) error

type DependencyProvider struct {
	factories map[reflect.Type]func(ctx context.Context, tag string) (any, Closer, error)
}

var (
	ErrInvalidDependencyFactory = errors.New("dependency factories must match the type di.DependencyFactory")
	ErrNotRegistered            = errors.New("dependency not registered")
	ErrFillParameters           = errors.New("invalid fill parameters")
)

// var defaultProvider = NewDependencyProvider()

func NewDependencyProvider() *DependencyProvider {
	return &DependencyProvider{
		factories: map[reflect.Type]func(ctx context.Context, tag string) (any, Closer, error){},
	}
}

func getType[T any]() reflect.Type {
	var v *T
	return reflect.TypeOf(v).Elem()
}
