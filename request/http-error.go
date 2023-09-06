package request

import "net/http"

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
	return e.err.Error()
}

func (e *HTTPError) Respond(w http.ResponseWriter, r *http.Request) error {
	return errorResponse(e.err, e.status, r).Respond(w, r)
}
