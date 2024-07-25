package openapidoc

import (
	"context"

	"github.com/go-openapi/spec"
)

type Operationer interface {
	Operation(ctx context.Context) (*spec.Operation, error)
}
type OperationMiddleware interface {
	OperationMiddleware(*spec.Operation) *spec.Operation
}
type Pathser interface {
	Paths(ctx context.Context, basePath string) (*spec.Paths, error)
}
type APIDocer interface {
	APIDoc(context.Context) (*spec.Swagger, error)
}
