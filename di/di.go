package di

import (
	"context"
	"errors"
	"reflect"
)

type DependencyProvider struct {
	factories map[reflect.Type]func(ctx context.Context, tag string) (any, error)
}

type DependencyFactory[T any] func(ctx context.Context, tag string) (T, error)

var (
	contextType            = getType[context.Context]()
	stringType             = getType[string]()
	errorType              = getType[error]()
	dependencyProviderType = getType[*DependencyProvider]()
)

var (
	ErrInvalidDependencyFactory = errors.New("dependency factories must match the type di.DependencyFactory")
	ErrNotRegistered            = errors.New("dependency not registered")
	ErrFillParameters           = errors.New("invalid fill parameters")
)

// var defaultProvider = NewDependencyProvider()

func NewDependencyProvider() *DependencyProvider {
	return &DependencyProvider{
		factories: map[reflect.Type]func(ctx context.Context, tag string) (any, error){},
	}
}

func getType[T any]() reflect.Type {
	var v *T
	return reflect.TypeOf(v).Elem()
}
