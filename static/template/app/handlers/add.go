package handlers

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/static/template/app/events"
)

type AddRequest struct {
	A   float64 `query:"a" validate:"required"`
	B   float64 `query:"b" validate:"required"`
	Ctx context.Context
}
type AddResponse struct {
	Sum float64 `json:"sum"`
}

var Add = request.Handler(func(r *AddRequest) (*AddResponse, error) {
	err := kernel.Dispatch(r.Ctx, &events.LogEvent{Message: fmt.Sprintf("add handler run with %f and %f", r.A, r.B)})
	if err != nil {
		return nil, err
	}
	return &AddResponse{
		Sum: r.A + r.B,
	}, nil
})
