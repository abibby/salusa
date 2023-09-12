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
