package request

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	data    any
	status  int
	headers map[string]string
}

var _ Responder = &JSONResponse{}

func NewJSONResponse(data any) *JSONResponse {
	return &JSONResponse{
		data:    data,
		headers: map[string]string{},
	}
}

func (r *JSONResponse) SetStatus(status int) *JSONResponse {
	r.status = status
	return r
}

func (r *JSONResponse) AddHeader(key, value string) *JSONResponse {
	r.headers[key] = value
	return r
}

func (r *JSONResponse) Respond(w http.ResponseWriter, _ *http.Request) error {
	if r.status != 0 {
		w.WriteHeader(r.status)
	}
	w.Header().Set("Content-Type", "application/json")
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	return e.Encode(r.data)
}
