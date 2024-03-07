package request

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abibby/salusa/router"
	"github.com/jmoiron/sqlx"
)

func HandleErrors(handlers ...func(err error)) router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					if err, ok := err.(error); ok {
						for _, handler := range handlers {
							handler(err)
						}
						if responder, ok := err.(Responder); ok {
							respond(w, r, responder)
						} else {
							respond(w, r, errorResponse(err, http.StatusInternalServerError, r))
						}
					} else {
						NewHTTPError(fmt.Errorf("internal server error"), 500).Respond(w, r)
					}
					return
				}

			}()

			next.ServeHTTP(w, r)
		})
	}
}

type txWrapper struct {
	tx *sqlx.Tx
}

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

func TransactionMiddleware() router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapper := &txWrapper{}

			defer func() {
				if r := recover(); r != nil {
					if wrapper.tx != nil {
						wrapper.tx.Rollback()
					}
					panic(r)
				}
			}()
			rw := &ResponseWriter{ResponseWriter: w, Status: 200}
			ctx := r.Context()
			ctx = context.WithValue(ctx, requestKey, r)
			ctx = context.WithValue(ctx, responseKey, w)
			ctx = context.WithValue(ctx, txKey, wrapper)
			next.ServeHTTP(rw, r.WithContext(ctx))

			if wrapper.tx != nil {
				tx := wrapper.tx
				if rw.OK() {
					err := tx.Commit()
					if err != nil {
						panic(err)
					}
				} else {
					err := tx.Rollback()
					if err != nil {
						panic(err)
					}
				}
			}
		})
	}
}
