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
					respond(w, r, responder)
				} else {
					respond(w, r, errorResponse(err, http.StatusInternalServerError, r))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
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
	return w.Status >= 200 && w.Status < 400
}
