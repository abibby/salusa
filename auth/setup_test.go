package auth_test

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/auth"
	"github.com/abibby/salusa/database/dbtest"
	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/database/migrate"
	"github.com/abibby/salusa/database/model"
	"github.com/jmoiron/sqlx"
)

type AutoIncrementUser struct {
	model.BaseModel
	ID           int    `db:"id,autoincrement,primary"`
	Username     string `db:"username"`
	PasswordHash []byte `db:"password"`
}

var _ auth.User = (*AutoIncrementUser)(nil)

func (u *AutoIncrementUser) GetID() string {
	return fmt.Sprint(u.ID)
}
func (u *AutoIncrementUser) GetPasswordHash() []byte {
	return u.PasswordHash
}
func (u *AutoIncrementUser) SetPasswordHash(p []byte) {
	u.PasswordHash = p
}
func (u *AutoIncrementUser) SaltedPassword(password string) []byte {
	return []byte(fmt.Sprintf("%d|%s", u.ID, password))
}
func (u *AutoIncrementUser) UsernameColumns() []string {
	return []string{"username"}
}

var runner = dbtest.NewRunner(func() (*sqlx.DB, error) {
	sqlite.UseSQLite()
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	err = migrate.RunModelCreate(ctx, db, &auth.UsernameUser{}, &auth.EmailVerifiedUser{}, &AutoIncrementUser{})
	if err != nil {
		return nil, err
	}
	return db, nil
})

var Run = runner.Run
