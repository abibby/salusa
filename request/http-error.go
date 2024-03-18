package request

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	err    error
	status int
}

var _ error = &HTTPError{}
var _ Responder = &HTTPError{}

func NewHTTPError(err error, status int) *HTTPError {
	return &HTTPError{
		err:    err,
		status: status,
	}
}
func (e *HTTPError) Error() string {
	return fmt.Sprintf("http %d: %v", e.status, e.err)
}
func (e *HTTPError) Unwrap() error {
	return e.err
}

func (e *HTTPError) Respond(w http.ResponseWriter, r *http.Request) error {
	return errorResponse(e.err, e.status, r).Respond(w, r)
}
