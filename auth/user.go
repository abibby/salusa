package auth

import (
	"github.com/abibby/salusa/database/model"
	"github.com/google/uuid"
)

type User interface {
	model.Model
	GetID() string
	GetUsername() string
	GetPasswordHash() []byte
	SetPasswordHash([]byte)
	SaltedPassword(password string) []byte
}

type BaseUser struct {
	model.BaseModel

	ID           uuid.UUID `json:"id"       db:"id,primary"`
	Username     string    `json:"username" db:"username"`
	PasswordHash []byte    `json:"-"        db:"password"`
}

var _ User = (*BaseUser)(nil)

func (u *BaseUser) GetID() string {
	return u.ID.String()
}
func (u *BaseUser) GetUsername() string {
	return u.Username
}
func (u *BaseUser) GetPasswordHash() []byte {
	return u.PasswordHash
}
func (u *BaseUser) SetPasswordHash(b []byte) {
	u.PasswordHash = b
}
func (u *BaseUser) SaltedPassword(password string) []byte {
	return []byte(u.ID.String() + password)
}

type EmailVerified interface {
	GetEmail() string
	SetValidationToken(string)
}

type MustVerifyEmail struct {
	Email           string `json:"email" db:"email"`
	Validated       bool   `json:"-"     db:"validated"`
	ValidationToken string `json:"-"     db:"validation_code"`
}

var _ EmailVerified = (*MustVerifyEmail)(nil)

func (v *MustVerifyEmail) GetEmail() string {
	return v.Email
}
func (v *MustVerifyEmail) SetValidationToken(t string) {
	v.ValidationToken = t
}
