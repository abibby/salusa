package auth

import (
	"fmt"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/golang-jwt/jwt/v4"
)

var ErrInvalidToken = fmt.Errorf("invalid token")

type ClaimType string

var (
	TypeAccess  = ClaimType("access")
	TypeRefresh = ClaimType("refresh")
)

type Claims struct {
	jwt.RegisteredClaims
	Type ClaimType `json:"type,omitempty"`
}

func ParseOf[T jwt.Claims](token string) (T, error) {
	claims := helpers.NewOf[T]()
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("expected HMAC received %v: %w", token.Header["alg"], ErrUnexpectedAlgorithm)
		}

		return appKey, nil
	})
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to parse JWT: %w", err)
	}
	if !t.Valid {
		var zero T
		return zero, ErrInvalidToken
	}
	return claims, nil
}
func Parse(token string) (*Claims, error) {
	return ParseOf[*Claims](token)
}
