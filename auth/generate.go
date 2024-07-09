package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(appKey)
}
