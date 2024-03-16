package auth

import (
	"github.com/abibby/salusa/database/model"
	"github.com/google/uuid"
)

type User interface {
	model.Model
	GetID() string
	GetUsername() string
	SetUsername(string)
	GetPasswordHash() []byte
	SetPasswordHash([]byte)
	SaltedPassword(password string) []byte
}

type EmailVerified interface {
	GetEmail() string
	SetValidationToken(string)
}

type BaseUser struct {
	model.BaseModel

	ID           uuid.UUID `json:"id"       db:"id,primary"`
	Username     string    `json:"username" db:"username"`
	PasswordHash []byte    `json:"-"        db:"password"`
}

var _ User = (*BaseUser)(nil)

func NewBaseUser() *BaseUser {
	return &BaseUser{
		ID: uuid.New(),
	}
}

func (u *BaseUser) GetID() string {
	return u.ID.String()
}
func (u *BaseUser) GetUsername() string {
	return u.Username
}
func (u *BaseUser) SetUsername(username string) {
	u.Username = username
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

type EmailVerifiedUser struct {
	model.BaseModel

	ID              uuid.UUID `json:"id"    db:"id,primary"`
	Email           string    `json:"email" db:"email"`
	PasswordHash    []byte    `json:"-"     db:"password"`
	Validated       bool      `json:"-"     db:"validated"`
	ValidationToken string    `json:"-"     db:"validation_code"`
}

var _ EmailVerified = (*EmailVerifiedUser)(nil)
var _ User = (*EmailVerifiedUser)(nil)

func NewEmailVerifiedUser() *EmailVerifiedUser {
	return &EmailVerifiedUser{
		ID: uuid.New(),
	}
}

func (u *EmailVerifiedUser) GetID() string {
	return u.ID.String()
}
func (u *EmailVerifiedUser) GetUsername() string {
	return u.Email
}
func (u *EmailVerifiedUser) SetUsername(username string) {
	u.Email = username
}
func (u *EmailVerifiedUser) GetPasswordHash() []byte {
	return u.PasswordHash
}
func (u *EmailVerifiedUser) SetPasswordHash(b []byte) {
	u.PasswordHash = b
}
func (u *EmailVerifiedUser) SaltedPassword(password string) []byte {
	return []byte(u.ID.String() + password)
}

func (v *EmailVerifiedUser) GetEmail() string {
	return v.Email
}
func (v *EmailVerifiedUser) SetValidationToken(t string) {
	v.ValidationToken = t
}
