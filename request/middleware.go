package request

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/abibby/salusa/router"
)

type errorContextKey struct{}

type errorsContainer struct {
	errors []error
}

func (e *errorsContainer) add(err error) {
	e.errors = append(e.errors, err)
}
func (w *errorsContainer) joinedError() error {
	if len(w.errors) == 0 {
		return nil
	}
	if len(w.errors) == 1 {
		return w.errors[0]
	}
	return errors.Join(w.errors...)
}

func addError(r *http.Request, err error) {
	e := r.Context().Value(errorContextKey{})
	if e == nil {
		return
	}
	e.(*errorsContainer).add(err)
}

func HandleErrors(handlers ...func(err error)) router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				e := recover()
				if e == nil {
					return
				}
				var err error
				if e, ok := e.(error); ok {
					err = e
				} else {
					err = fmt.Errorf("internal server error")
				}
				for _, handler := range handlers {
					handler(err)
				}

				err = fmt.Errorf("%w\n%s", err, debug.Stack())
				responder, ok := getResponder(err)
				if !ok {
					responder = NewHTTPError(err, http.StatusInternalServerError)
				}
				err = responder.Respond(w, r)
				if err != nil {
					log.Print(err)
				}
			}()

			e := &errorsContainer{errors: []error{}}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), errorContextKey{}, e)))

			err := e.joinedError()
			if err != nil {
				for _, handler := range handlers {
					handler(err)
				}
			}
		})
	}
}
