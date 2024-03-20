package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
	"github.com/golang-jwt/jwt/v4"
)

type contextKey uint8

const (
	jwtClaims contextKey = iota
)

var appKey []byte

func SetAppKey(key []byte) {
	appKey = key
}

func GetClaims(ctx context.Context) (jwt.MapClaims, bool) {
	iClaims := ctx.Value(jwtClaims)
	claims, ok := iClaims.(jwt.MapClaims)
	return claims, ok
}

func UserIDFactory[T any](cb func(claims jwt.MapClaims) (T, bool)) func(ctx context.Context) (T, bool) {
	return func(ctx context.Context) (T, bool) {
		var zero T
		claims, ok := GetClaims(ctx)
		if !ok {
			return zero, false
		}

		return cb(claims)
	}
}

func setClaims(r *http.Request, claims *Claims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), jwtClaims, claims))
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
		_, ok := GetClaims(r.Context())
		if !ok {
			respond(request.NewHTTPError(fmt.Errorf("unauthorized"), http.StatusUnauthorized), w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func HasClaim(key string, value any) router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok {
				respond(request.NewHTTPError(fmt.Errorf("unauthorized"), http.StatusUnauthorized), w, r)
				return
			}
			claim, ok := claims[key]
			if !ok {
				respond(request.NewHTTPError(fmt.Errorf("unauthorized"), http.StatusUnauthorized), w, r)
				return
			}
			if claim != value {
				respond(request.NewHTTPError(fmt.Errorf("unauthorized"), http.StatusUnauthorized), w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func respond(resp request.Responder, w http.ResponseWriter, r *http.Request) {
	err := resp.Respond(w, r)
	if err != nil {
		log.Print(err)
	}
}
