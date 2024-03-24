package auth

import (
	"github.com/abibby/salusa/database/model"
	"github.com/google/uuid"
)

type User interface {
	model.Model
	GetID() string
	GetPasswordHash() []byte
	SetPasswordHash([]byte)
	SaltedPassword(password string) []byte
	UsernameColumns() []string
}

type EmailVerified interface {
	GetEmail() string
	SetLookupToken(string)
	IsVerified() bool
	SetVerified(bool)
	LookupTokenColumn() string
}

type UsernameUser struct {
	model.BaseModel

	ID           uuid.UUID `json:"id"       db:"id,primary"`
	Username     string    `json:"username" db:"username,unique"`
	PasswordHash []byte    `json:"-"        db:"password"`
}

var _ User = (*UsernameUser)(nil)

type UsernameUserCreateRequest struct {
	Username string `json:"username" validate:"required"`
}

func NewUsernameUser(request *UsernameUserCreateRequest) *UsernameUser {
	return &UsernameUser{
		ID:           uuid.New(),
		Username:     request.Username,
		PasswordHash: []byte{},
	}
}

func (u *UsernameUser) GetID() string {
	return u.ID.String()
}
func (u *UsernameUser) GetUsername() string {
	return u.Username
}
func (u *UsernameUser) UsernameColumns() []string {
	return []string{"username"}
}
func (u *UsernameUser) GetPasswordHash() []byte {
	return u.PasswordHash
}
func (u *UsernameUser) SetPasswordHash(b []byte) {
	u.PasswordHash = b
}
func (u *UsernameUser) SaltedPassword(password string) []byte {
	return []byte(u.ID.String() + password)
}
func (u *UsernameUser) PasswordColumn() string {
	return "password"
}

type EmailVerifiedUser struct {
	model.BaseModel

	ID           uuid.UUID `json:"id"    db:"id,primary"`
	Email        string    `json:"email" db:"email,unique"`
	PasswordHash []byte    `json:"-"     db:"password"`
	Verified     bool      `json:"-"     db:"validated"`
	LookupToken  string    `json:"-"     db:"lookup_token"`
}

var _ EmailVerified = (*EmailVerifiedUser)(nil)
var _ User = (*EmailVerifiedUser)(nil)

type EmailVerifiedUserCreateRequest struct {
	Email string `json:"username"  validate:"required|email"`
}

func NewEmailVerifiedUser(request *EmailVerifiedUserCreateRequest) *EmailVerifiedUser {
	return &EmailVerifiedUser{
		ID:           uuid.New(),
		Email:        request.Email,
		PasswordHash: []byte{},
	}
}

func (u *EmailVerifiedUser) GetID() string {
	return u.ID.String()
}
func (u *EmailVerifiedUser) GetUsername() string {
	return u.Email
}
func (u *EmailVerifiedUser) UsernameColumns() []string {
	return []string{"email"}
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
func (u *EmailVerifiedUser) LookupTokenColumn() string {
	return "lookup_token"
}

func (v *EmailVerifiedUser) GetEmail() string {
	return v.Email
}
func (v *EmailVerifiedUser) SetLookupToken(t string) {
	v.LookupToken = t
}
func (v *EmailVerifiedUser) IsVerified() bool {
	return v.Verified
}
func (v *EmailVerifiedUser) SetVerified(verified bool) {
	v.Verified = verified
}
