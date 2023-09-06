package request

import (
	_ "embed"
	"io"
	"net/http"
	"strings"
	"text/template"
)

//go:embed error.html
var errorTemplate string

type ErrResponse struct {
	Error      string          `json:"error"`
	Status     int             `json:"status"`
	StatusText string          `json:"-"`
	Fields     ValidationError `json:"fields,omitempty"`
}

func ErrorResponse(err error, status int, r *http.Request) Responder {
	response := ErrResponse{
		Error:      err.Error(),
		Status:     status,
		StatusText: http.StatusText(status),
	}
	if validationErr, ok := err.(ValidationError); ok {
		response.Fields = validationErr
	}
	if strings.HasPrefix(r.Header.Get("Accept"), "text/html") {
		t, err := template.New("error").Parse(errorTemplate)
		if err != nil {
			panic(err)
		}
		reader, writer := io.Pipe()
		go func() {
			err := t.Execute(writer, response)
			if err != nil {
				panic(err)
			}
			err = writer.Close()
			if err != nil {
				panic(err)
			}
		}()
		return NewResponse(reader).SetStatus(status).AddHeader("Content-Type", "text/html")
	}
	return NewJSONResponse(response).SetStatus(status)
}
