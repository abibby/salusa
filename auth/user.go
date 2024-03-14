package auth

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
)

type User struct {
	model.BaseModel

	ID           int    `json:"id"       db:"id,primary,autoincrement"`
	Username     string `json:"username" db:"username"`
	PasswordHash []byte `json:"-"        db:"password"`
}

func UserQuery(ctx context.Context) *builder.Builder[*User] {
	return builder.From[*User]().WithContext(ctx)
}

func (u *User) SaltedPassword(password string) []byte {
	return []byte(fmt.Sprintf("%d%s", u.ID, password))
}
