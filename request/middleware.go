package request

import (
	"fmt"
	"log"
	"net/http"

	"github.com/abibby/salusa/router"
)

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

				responder, ok := getResponder(err)
				if !ok {
					responder = errorResponse(err, http.StatusInternalServerError, r)
				}
				err = responder.Respond(w, r)
				if err != nil {
					log.Print(err)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
