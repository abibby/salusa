package modeldi

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/request"
	"github.com/jmoiron/sqlx"
)

func Register[T model.Model]() {
	di.Register(func(ctx context.Context, tag string) (T, error) {
		var zero T
		if tag == "" {
			return zero, fmt.Errorf("no tag")
		}
		req, err := di.Resolve[*http.Request](ctx)
		if err != nil {
			return zero, fmt.Errorf("count not find request: %w", err)
		}
		db, err := di.Resolve[*sqlx.DB](ctx)
		if err != nil {
			return zero, fmt.Errorf("count not find db: %w", err)
		}

		v := req.URL.Query().Get(tag)

		u, err := builder.From[T]().WithContext(ctx).Find(db, v)
		if err != nil {
			return zero, err
		}
		if reflect.ValueOf(u).IsZero() {
			return zero, request.NewHTTPError(fmt.Errorf("404 not found"), 404)
		}
		return u, nil
	})
}
