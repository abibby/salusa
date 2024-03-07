package request

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/abibby/salusa/di"
	"github.com/jmoiron/sqlx"
)

type contextKey struct{}

var (
	txKey       = contextKey{}
	requestKey  = contextKey{}
	responseKey = contextKey{}
)

func Init(context.Context) error {
	di.Register(func(ctx context.Context, tag string) *sqlx.Tx {
		wrapper, ok := ctx.Value(txKey).(*txWrapper)
		if !ok {
			log.Print("no transaction wrapper")
			return nil
		}

		if wrapper.tx == nil {
			db, ok := di.Resolve[*sqlx.DB](ctx)
			if !ok {
				return nil
			}
			tx, err := db.BeginTxx(ctx, &sql.TxOptions{
				ReadOnly: strings.ToLower(tag) == "r",
			})
			if err != nil {
				return nil
			}
			wrapper.tx = tx
		}
		return wrapper.tx
	})
	di.Register(func(ctx context.Context, tag string) *http.Request {
		req, ok := ctx.Value(requestKey).(*http.Request)
		if !ok {
			return nil
		}
		return req
	})
	di.Register(func(ctx context.Context, tag string) http.ResponseWriter {
		resp, ok := ctx.Value(responseKey).(http.ResponseWriter)
		if !ok {
			return nil
		}
		return resp
	})
	return nil
}
