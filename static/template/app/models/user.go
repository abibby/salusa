package models

import (
	"context"
	"crypto/sha512"
	"log"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/hooks"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/database/model/modeldi"
	"github.com/abibby/salusa/static/template/app"
)

//go:generate spice generate:migration
type User struct {
	model.BaseModel

	ID           int    `json:"id"       db:"id,primary,autoincrement"`
	Username     string `json:"username" db:"username"`
	Password     []byte `json:"-"        db:"-"`
	PasswordHash []byte `json:"-"        db:"password"`
}

var _ hooks.BeforeSaver = (*User)(nil)

func init() {
	app.Kernel.Register(modeldi.Register[*User])
}

func UserQuery(ctx context.Context) *builder.Builder[*User] {
	return builder.From[*User]().WithContext(ctx)
}

func (u *User) BeforeSave(ctx context.Context, db database.DB) error {
	if u.Password != nil {
		log.Print("THIS IS JUST FOR AN EXAMPLE. REPLACE THIS")
		h := sha512.Sum512_256(u.Password)
		u.PasswordHash = h[:]
	}
	return nil
}
