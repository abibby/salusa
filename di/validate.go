package di

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/validate"
	"github.com/dominikbraun/graph"
)

var (
	ErrDependancyCycle = errors.New("dependancy cycle")
)

type DIValidator struct {
	dp  *DependencyProvider
	typ reflect.Type
}

var _ validate.Validator = (*DIValidator)(nil)

func (v *DIValidator) Validate(ctx context.Context) error {
	errs := []error{}
	t := v.typ
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	for _, sf := range helpers.GetFields(t) {
		if !sf.IsExported() {
			continue
		}

		_, ok := sf.Tag.Lookup("inject")
		if !ok {
			continue
		}

		switch sf.Type {
		case contextType, dependencyProviderType:
			continue
		}

		_, ok = v.dp.factories[sf.Type]
		if !ok {
			errs = append(errs, fmt.Errorf("missing dependancy %s on %s.%s", sf.Type, v.typ, sf.Name))
		}
	}

	return errors.Join(errs...)
}

func Validator(ctx context.Context, rootType reflect.Type) *DIValidator {
	return GetDependencyProvider(ctx).Validator(rootType)
}
func (dp *DependencyProvider) Validator(rootType reflect.Type) *DIValidator {
	return &DIValidator{
		dp:  dp,
		typ: rootType,
	}
}

func (dp *DependencyProvider) Validate(ctx context.Context) error {
	errs := []error{}

	if err := dp.validateCycles(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func typeHash(t reflect.Type) reflect.Type {
	return t
}
func (dp *DependencyProvider) validateCycles() error {
	g := graph.New(typeHash, graph.Directed(), graph.PreventCycles())

	var err error
	for typ := range dp.factories {
		err = g.AddVertex(typ)
		if err != nil {
			panic(err)
		}
	}
	for typ, factory := range dp.factories {
		depends, ok := factory.(DependsOner)
		if !ok {
			continue
		}
		deps := depends.DependsOn()
		for _, dep := range deps {
			err = g.AddEdge(typ, dep)
			if errors.Is(err, graph.ErrEdgeCreatesCycle) {
				return newCycleError(g, dep, typ)
			} else if err != nil {
				panic(err)
			}
		}
	}
	return nil
}
func newCycleError[K comparable, T any](g graph.Graph[K, T], source, target K) error {
	path, err := graph.ShortestPath(g, source, target)
	if err != nil {
		panic(err)
	}

	strPath := []byte{}
	for _, t := range path {
		strPath = fmt.Appendf(strPath, "%v -> ", t)
	}
	strPath = fmt.Appendf(strPath, "%v", path[0])

	return fmt.Errorf("%w: %s", ErrDependancyCycle, strPath)

}
