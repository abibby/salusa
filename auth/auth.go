package auth

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

type contextKey uint8

const (
	claimKey contextKey = iota
)

var (
	Err401Unauthorized = request.NewHTTPError(errors.New("unauthorized"), http.StatusUnauthorized)
)

var appKey []byte

func SetAppKey(key []byte) {
	appKey = key
}

func GetClaims(r *http.Request) (*Claims, bool) {
	return GetClaimsCtx(r.Context())
}
func GetClaimsCtx(ctx context.Context) (*Claims, bool) {
	iClaims := ctx.Value(claimKey)
	if iClaims == nil {
		return nil, false
	}
	claims, ok := iClaims.(*Claims)
	return claims, ok
}

func setClaims(r *http.Request, claims *Claims) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), claimKey, claims))
}

func AttachUser() router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := authenticate(r)
			if err == nil {
				next.ServeHTTP(w, setClaims(r, claims))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func LoggedIn() router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := GetClaims(r)
			if !ok {
				respond(Err401Unauthorized, w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func HasClaim(validate func(c *Claims) bool) router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r)
			if !ok {
				respond(Err401Unauthorized, w, r)
				return
			}

			if !validate(claims) {
				respond(Err401Unauthorized, w, r)
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
