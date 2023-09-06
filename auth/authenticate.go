package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrInvalidAuthorizationHeader = fmt.Errorf("missing or invalid Authorization header")
	ErrUnexpectedAlgorithm        = fmt.Errorf("unexpected algorithm")
	ErrInvalidToken               = fmt.Errorf("invalid token")
)

func authenticate(r *http.Request) (jwt.MapClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, ErrInvalidAuthorizationHeader
	}
	tokenStr := authHeader[7:]

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("expected HMAC received %v: %w", token.Header["alg"], ErrUnexpectedAlgorithm)
		}

		return appKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
