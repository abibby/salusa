package request

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/helpers"
	"github.com/abibby/salusa/router"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

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
	return w.Status >= 200 && w.Status < 300
}

type contextKey int

const (
	dbKey contextKey = iota
	txKey
)

var (
	ErrNoTx                = fmt.Errorf("no transaction set")
	ErrNotModel            = fmt.Errorf("propery must be of type models.Model")
	ErrCompositePrimaryKey = fmt.Errorf("can't use models with composite primary keys")
)

func WithDB(db *sqlx.DB) router.MiddlewareFunc {
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
			ctx = context.WithValue(ctx, dbKey, db)
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

func SetTestTx(r *http.Request, tx *sqlx.Tx) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), txKey, &txWrapper{tx: tx}))
}

func UseDB(r *http.Request) *sqlx.DB {
	if db := r.Context().Value(dbKey); db != nil {
		return db.(*sqlx.DB)
	}
	return nil
}

func UseTx(r *http.Request) *sqlx.Tx {
	if tx := r.Context().Value(txKey); tx != nil {
		wrapper := tx.(*txWrapper)
		if wrapper.tx == nil {
			db := UseDB(r)
			tx, err := db.BeginTxx(r.Context(), &sql.TxOptions{})
			if err != nil {
				panic(err)
			}
			wrapper.tx = tx
		}
		return wrapper.tx
	}
	return nil
}

func di(v any, r *http.Request) error {
	vars := mux.Vars(r)
	r.Context()
	return helpers.EachField(reflect.ValueOf(v), func(sf reflect.StructField, fv reflect.Value) error {
		tag, ok := sf.Tag.Lookup("di")
		if !ok {
			return nil
		}
		id, ok := vars[tag]
		if !ok {
			return nil
		}
		tx := UseTx(r)
		if tx == nil {
			return ErrNoTx
		}

		f, ok := fv.Interface().(model.Model)
		if !ok {
			return ErrNotModel
		}
		pKey := helpers.PrimaryKey(f)
		if len(pKey) != 1 {
			return ErrCompositePrimaryKey
		}
		q := builder.New[model.Model]().
			From(helpers.GetTable(f)).
			Where(pKey[0], "=", id).
			WithContext(r.Context())

		if fv.Kind() != reflect.Pointer {
			return fmt.Errorf("expected pointer")
		}
		f = reflect.New(fv.Type().Elem()).Interface().(model.Model)
		err := q.LoadOne(tx, f)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		} else if err != nil {
			return err
		}
		newValue := reflect.ValueOf(f)
		if !newValue.IsZero() {
			fv.Set(newValue)
		}
		return nil
	})
}
