package databasedi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abibby/salusa/router"
	"github.com/jmoiron/sqlx"
)

type txWrapper struct {
	tx *sqlx.Tx
}

type contextKey uint8

const (
	txKey contextKey = iota
	dbKey
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
	ErrNoTx                = fmt.Errorf("no transaction set")
	ErrNotModel            = fmt.Errorf("propery must be of type models.Model")
	ErrCompositePrimaryKey = fmt.Errorf("can't use models with composite primary keys")
)

func Middleware() router.MiddlewareFunc {
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
