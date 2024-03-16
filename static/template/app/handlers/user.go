package handlers

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/static/template/app/events"
	"github.com/abibby/salusa/static/template/app/models"
	"github.com/jmoiron/sqlx"
)

type ListUserRequest struct {
	Tx  *sqlx.Tx        `inject:""`
	Ctx context.Context `inject:""`
}
type ListUserResponse struct {
	Users []*models.User `json:"users"`
}

var UserList = request.Handler(func(r *ListUserRequest) (*ListUserResponse, error) {
	users, err := models.UserQuery(r.Ctx).Get(r.Tx)
	if err != nil {
		return nil, err
	}
	return &ListUserResponse{
		Users: users,
	}, nil
})

type GetUserRequest struct {
	User *models.User `inject:"id"`
}
type GetUserResponse struct {
	User *models.User `json:"user"`
}

var UserGet = request.Handler(func(r *GetUserRequest) (*GetUserResponse, error) {
	return &GetUserResponse{
		User: r.User,
	}, nil
})

type CreateUserRequest struct {
	Username string          `json:"username" validate:"required"`
	Password []byte          `json:"password" validate:"required"`
	Tx       *sqlx.Tx        `inject:""`
	Ctx      context.Context `inject:""`
	Queue    event.Queue     `inject:""`
}
type CreateUserResponse struct {
	User *models.User `json:"user"`
}

var UserCreate = request.Handler(func(r *CreateUserRequest) (*CreateUserResponse, error) {
	u := &models.User{
		// Username: r.Username,
		// Password: r.Password,
	}

	err := model.SaveContext(r.Ctx, r.Tx, u)
	if err != nil {
		return nil, err
	}

	err = r.Queue.Push(&events.LogEvent{Message: fmt.Sprintf("create user with username %s", r.Username)})
	if err != nil {
		return nil, err
	}

	return &CreateUserResponse{
		User: u,
	}, nil
})
