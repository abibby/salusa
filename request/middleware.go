package request

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/abibby/salusa/router"
)

type errorContextType uint8

const errorContextKey = errorContextType(iota)

type errorsContainer struct {
	errors []error
}

func (e *errorsContainer) add(err error) {
	e.errors = append(e.errors, err)
}

func addError(r *http.Request, err error) {
	e := r.Context().Value(errorContextKey)
	if e == nil {
		return
	}
	e.(*errorsContainer).add(err)
}

func HandleErrors(handlers ...func(ctx context.Context, err error)) router.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			errContainer := &errorsContainer{errors: []error{}}

			defer func() {
				recoverErr := toError(recover())
				if recoverErr != nil {
					errContainer.add(recoverErr)
				}

				for _, err := range errContainer.errors {
					for _, handler := range handlers {
						handler(r.Context(), err)
					}
				}

				if recoverErr == nil {
					return
				}

				responder, ok := getResponder(recoverErr)
				if !ok {
					responder = NewHTTPError(recoverErr, http.StatusInternalServerError)
				}
				if httpErr, ok := responder.(*HTTPError); ok {
					httpErr.WithStack()
				}
				recoverErr = responder.Respond(w, r)
				if recoverErr != nil {
					log.Print(recoverErr)
				}
			}()

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), errorContextKey, errContainer)))
		})
	}
}

func toError(e any) error {
	if e == nil {
		return nil
	}
	switch e := e.(type) {
	case error:
		return e
	case string:
		return errors.New(e)
	default:
		return NewDefaultHTTPError(http.StatusInternalServerError)
	}
}
