package handlers

import (
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/static/template/app/events"
)

type AddRequest struct {
	A float64 `query:"a" validate:"required"`
	B float64 `query:"b" validate:"required"`
}
type AddResponse struct {
	Sum float64 `json:"sum"`
}

var Add = request.Handler(func(r *AddRequest) (*AddResponse, error) {
	kernel.Dispatch(&events.LogEvent{Message: "add handler run"})
	return &AddResponse{
		Sum: r.A + r.B,
	}, nil
})
