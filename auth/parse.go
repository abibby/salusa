package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func Parse(token string) (jwt.MapClaims, error) {

	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("expected HMAC received %v: %w", token.Header["alg"], ErrUnexpectedAlgorithm)
		}

		return appKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}
	if !t.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
