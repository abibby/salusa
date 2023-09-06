package requests

import "net/http"

type HTTPError struct {
	message string
	status  int
}

var _ error = &HTTPError{}
var _ Responder = &HTTPError{}

func (e *HTTPError) Error() string {
	return e.message
}

func (e *HTTPError) Respond(w http.ResponseWriter) error {
	w.WriteHeader(e.status)
	w.Write([]byte(e.message))
	return nil
}
