package handlers

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/kernel"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/jmoiron/sqlx"
)

type GetUserRequest struct {
	User *models.User `inject:"user_id"`
}
type GetUaerResponse struct {
	User *models.User `json:"user"`
}

var GetUser = request.Handler(func(r *GetUserRequest) (*GetUaerResponse, error) {
	return &GetUaerResponse{
		User: r.User,
	}, nil
})

type CreateUserRequest struct {
	Username string          `json:"username" validate:"required"`
	Password []byte          `json:"password" validate:"required"`
	Tx       *sqlx.Tx        `inject:""`
	Ctx      context.Context `inject:""`
}
type CreateUaerResponse struct {
	User *models.User `json:"user"`
}

var CreateUser = request.Handler(func(r *CreateUserRequest) (*GetUaerResponse, error) {
	err := kernel.Dispatch(r.Ctx, &events.LogEvent{Message: fmt.Sprintf("create user with username %s", r.Username)})
	if err != nil {
		return nil, err
	}

	u := &models.User{
		Username: r.Username,
		Password: r.Password,
	}

	err = model.SaveContext(r.Ctx, r.Tx, u)
	if err != nil {
		return nil, err
	}

	return &GetUaerResponse{
		User: u,
	}, nil
})
