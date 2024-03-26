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
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type modelDeps struct {
	Request *http.Request `inject:""`
	DB      *sqlx.DB      `inject:""`
}

func Register[T model.Model](ctx context.Context) {
	di.RegisterWith(ctx, func(ctx context.Context, tag string, deps *modelDeps) (T, error) {
		var zero T

		v, ok := getValue(deps.Request, tag)
		if !ok {
			return zero, fmt.Errorf("could not fetch model %s not in request", tag)
		}

		u, err := builder.From[T]().WithContext(ctx).Find(deps.DB, v)
		if err != nil {
			return zero, err
		}
		if reflect.ValueOf(u).IsZero() {
			return zero, request.NewHTTPError(fmt.Errorf("404 not found"), 404)
		}
		return u, nil
	})
}

func getValue(r *http.Request, tag string) (string, bool) {
	if r.URL.Query().Has(tag) {
		return r.URL.Query().Get(tag), true
	}

	vars := mux.Vars(r)
	v, ok := vars[tag]
	if ok {
		return v, true
	}
	return "", false
}
