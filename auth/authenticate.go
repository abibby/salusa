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
	return Parse(tokenStr)
}
