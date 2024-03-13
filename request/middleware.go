package request

import (
	"fmt"
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

				if responder, ok := getResponder(err); ok {
					responder.Respond(w, r)
				} else {
					errorResponse(err, http.StatusInternalServerError, r).Respond(w, r)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
