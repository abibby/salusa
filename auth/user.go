package auth

import (
	"github.com/abibby/salusa/database/model"
	"github.com/google/uuid"
)

type User interface {
	model.Model
	GetID() string
	UsernameColumn() string
	GetUsername() string
	SetUsername(string)
	GetPasswordHash() []byte
	SetPasswordHash([]byte)
	SaltedPassword(password string) []byte
	PasswordColumn() string
}

type EmailVerified interface {
	GetEmail() string
	SetVerificationToken(string)
	IsVerified() bool
	SetVerified(bool)
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
func (u *BaseUser) UsernameColumn() string {
	return "username"
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
func (u *BaseUser) PasswordColumn() string {
	return "password"
}

type EmailVerifiedUser struct {
	model.BaseModel

	ID                uuid.UUID `json:"id"    db:"id,primary"`
	Email             string    `json:"email" db:"email"`
	PasswordHash      []byte    `json:"-"     db:"password"`
	Verified          bool      `json:"-"     db:"validated"`
	VerificationToken string    `json:"-"     db:"verification_token"`
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
func (u *EmailVerifiedUser) UsernameColumn() string {
	return "email"
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
func (u *EmailVerifiedUser) PasswordColumn() string {
	return "password"
}

func (v *EmailVerifiedUser) GetEmail() string {
	return v.Email
}
func (v *EmailVerifiedUser) SetVerificationToken(t string) {
	v.VerificationToken = t
}
func (v *EmailVerifiedUser) IsVerified() bool {
	return v.Verified
}
func (v *EmailVerifiedUser) SetVerified(verified bool) {
	v.Verified = verified
}
