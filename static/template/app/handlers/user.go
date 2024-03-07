package handlers

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/jmoiron/sqlx"
)

type UserRequest struct {
	ID  int             `query:"user_id" validate:"required"`
	TX  *sqlx.DB        `inject:""`
	Ctx context.Context `inject:""`
}
type UaerResponse struct {
	User *models.User `json:"user"`
}

var User = request.Handler(func(r *UserRequest) (*UaerResponse, error) {
	err := kernel.Dispatch(r.Ctx, &events.LogEvent{Message: fmt.Sprintf("fetch user wiht id %d", r.ID)})
	if err != nil {
		return nil, err
	}

	u, err := models.UserQuery(r.Ctx).Find(r.TX, r.ID)
	if err != nil {
		return nil, err
	}

	return &UaerResponse{
		User: u,
	}, nil
})
