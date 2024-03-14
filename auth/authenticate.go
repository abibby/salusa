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
	prefix := "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, prefix) {
		return nil, ErrInvalidAuthorizationHeader
	}
	tokenStr := authHeader[len(prefix):]
	return Parse(tokenStr)
}
