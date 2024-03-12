package databasedi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abibby/salusa/router"
)

type ResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) OK() bool {
	return w.Status >= 200 && w.Status < 400
}

var (
	ErrRequestFailed = fmt.Errorf("request failed")
)

func Middleware() router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx, cancel := context.WithCancelCause(ctx)

			defer func() {
				if r := recover(); r != nil {
					cancel(ErrRequestFailed)
					panic(r)
				}
			}()

			rw := &ResponseWriter{ResponseWriter: w, Status: 200}
			next.ServeHTTP(rw, r.WithContext(ctx))

			if !rw.OK() {
				cancel(ErrRequestFailed)
			} else {
				cancel(nil)
			}
		})
	}
}
