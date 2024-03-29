package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
)

func SetClaims(ctx context.Context, claims *Claims) context.Context {
	return setClaims(httptest.NewRequest(http.MethodGet, "/", http.NoBody).WithContext(ctx), claims).Context()
}
