package request

import (
	"bytes"
	"fmt"
	"html"
)

func Redirect(to string) *Response {
	b := &bytes.Buffer{}
	fmt.Fprintf(b, `<a href="%s">redirect</a>`, html.EscapeString(to))
	return NewResponse(b).
		SetStatus(301).
		AddHeader("Location", to)
}
