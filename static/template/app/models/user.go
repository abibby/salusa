package models

import (
	"context"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/jmoiron/sqlx"
)

//go:generate spice generate:migration
type User struct {
	model.BaseModel

	ID           int    `json:"id"       db:"id,primary,autoincrement"`
	Username     string `json:"username" db:"username"`
	Password     []byte `json:"-"        db:"-"`
	PasswordHash []byte `json:"-"        db:"password"`
}

func UserQuery(ctx context.Context) *builder.Builder[*User] {
	return builder.From[*User]().WithContext(ctx)
}

func (u *User) Save(tx *sqlx.Tx) error {
	return model.Save(tx, u)
}
