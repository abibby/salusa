package auth

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

var (
	ErrMissingAuthorizationHeader = fmt.Errorf("missing Authorization header")
	ErrInvalidAuthorizationHeader = fmt.Errorf("invalid Authorization header")
	ErrUnexpectedAlgorithm        = fmt.Errorf("unexpected algorithm")
	ErrNoAccessScope              = fmt.Errorf("no access scope")
)

func authenticate(r *http.Request) (*Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrMissingAuthorizationHeader
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return nil, ErrInvalidAuthorizationHeader
	}
	tokenStr := authHeader[len(prefix):]
	claims, err := Parse(tokenStr)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(claims.Scope, ScopeAccess) {
		return nil, ErrNoAccessScope
	}
	return claims, nil
}
