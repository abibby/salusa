package request

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type ValidationError map[string][]string

var _ error = ValidationError{}

func (e ValidationError) Error() string {
	// return "validation error"
	return fmt.Sprintf("validation error (%d)", len(e))
}
func (e ValidationError) HTMLError() string {
	items := make([]string, 0, len(e))
	for field, errors := range e {
		b := &bytes.Buffer{}
		fmt.Fprint(b, "<li>"+field+"<ul>")
		for _, err := range errors {
			fmt.Fprint(b, "<li>"+err+"</li>")
		}
		fmt.Fprint(b, "</ul></li>")
		items = append(items, b.String())
	}
	slices.Sort(items)
	return fmt.Sprintf(
		"<h2>Validation Error</h2><ul>%s</ul>",
		strings.Join(items, ""),
	)
}

func (e ValidationError) HasErrors() bool {
	return len(e) > 0
}

func (e ValidationError) AddError(key string, message string) {
	messages := e[key]
	if messages == nil {
		messages = []string{message}
	} else {
		messages = append(messages, message)
	}
	e[key] = messages
}
func (e ValidationError) Merge(vErr ValidationError) {
	for k, v := range vErr {
		errs := e[k]
		if errs == nil {
			errs = v
		} else {
			errs = append(errs, v...)
		}
		e[k] = errs
	}
}
