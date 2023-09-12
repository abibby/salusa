package controllers

import (
	"github.com/abibby/salusa/request"
)

type AddRequest struct {
	A float64 `query:"a" validate:"required"`
	B float64 `query:"b" validate:"required"`
}
type AddResponse struct {
	Sum float64 `json:"sum"`
}

var Add = request.Handler(func(r *AddRequest) (*AddResponse, error) {
	return &AddResponse{
		Sum: r.A + r.B,
	}, nil
})
