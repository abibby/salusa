package request

import (
	_ "embed"
	"io"
	"net/http"
	"strings"
	"text/template"
)

type HTMLError interface {
	HTMLError() string
}

//go:embed error.html
var errorTemplate string

type errResponse struct {
	Error      string          `json:"error"`
	Status     int             `json:"status"`
	StatusText string          `json:"-"`
	Fields     ValidationError `json:"fields,omitempty"`
}

func errorResponse(rootErr error, status int, r *http.Request) Responder {
	response := errResponse{
		Error:      rootErr.Error(),
		Status:     status,
		StatusText: http.StatusText(status),
	}
	if validationErr, ok := rootErr.(ValidationError); ok {
		response.Fields = validationErr
	}
	if strings.HasPrefix(r.Header.Get("Accept"), "text/html") {

		t, err := template.New("error").Parse(errorTemplate)
		if err != nil {
			panic(err)
		}
		if err, ok := rootErr.(HTMLError); ok {
			response.Error = err.HTMLError()
		} else {
			response.Error = "<p>" + rootErr.Error() + "</p>"
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
