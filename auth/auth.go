package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/golang-jwt/jwt/v4"
)

var appKey string

func SetAppKey(key string) {
	appKey = key
}

func Claims(ctx context.Context) (jwt.MapClaims, bool) {
	iClaims := ctx.Value("jwt-claims")
	claims, ok := iClaims.(jwt.MapClaims)
	return claims, ok
}

func UserIDFactory[T any](cb func(claims jwt.MapClaims) (T, bool)) func(ctx context.Context) (T, bool) {
	return func(ctx context.Context) (T, bool) {
		var zero T
		claims, ok := Claims(ctx)
		if !ok {
			return zero, false
		}

		return cb(claims)
	}
}

func setClaims(r *http.Request, claims jwt.MapClaims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "jwt-claims", claims))
}

func AttachUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := authenticate(r)
		if err == nil {
			next.ServeHTTP(w, setClaims(r, claims))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := Claims(r.Context())
		if !ok {
			request.ErrorResponse(fmt.Errorf("unauthorized"), http.StatusUnauthorized, r).Respond(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func HasClaim(key string, value any) router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := Claims(r.Context())
			if !ok {
				request.ErrorResponse(fmt.Errorf("unauthorized"), http.StatusUnauthorized, r).Respond(w)
				return
			}
			claim, ok := claims[key]
			if !ok {
				request.ErrorResponse(fmt.Errorf("unauthorized"), http.StatusUnauthorized, r).Respond(w)
				return
			}
			if claim != value {
				request.ErrorResponse(fmt.Errorf("unauthorized"), http.StatusUnauthorized, r).Respond(w)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
