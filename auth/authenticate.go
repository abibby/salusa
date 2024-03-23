package auth

import (
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrInvalidAuthorizationHeader = fmt.Errorf("missing or invalid Authorization header")
	ErrUnexpectedAlgorithm        = fmt.Errorf("unexpected algorithm")
)

func authenticate(r *http.Request) (*Claims, error) {
	prefix := "Bearer "
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, prefix) {
		return nil, ErrInvalidAuthorizationHeader
	}
	tokenStr := authHeader[len(prefix):]
	return Parse(tokenStr)
}
