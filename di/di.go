package di

import (
	"context"
	"errors"
	"reflect"
)

type DependencyProvider struct {
	factories map[reflect.Type]func(ctx context.Context, tag string) (any, error)
}

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
