package handlers

import (
	"context"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
	"github.com/go-openapi/spec"
)

type DocsRequest struct {
	Ctx context.Context `inject:""`
}

var Docs = request.Handler(func(r *DocsRequest) (*spec.Swagger, error) {
	k, err := di.Resolve[*kernel.Kernel](r.Ctx)
	if err != nil {
		return nil, err
	}
	return k.APIDoc(r.Ctx)
})
