package request

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

type HTMLError interface {
	HTMLError() string
}

//go:embed error.html
var errorTemplate string

type errResponse struct {
	Error      template.HTML   `json:"error"`
	Status     int             `json:"status"`
	StatusText string          `json:"-"`
	StackTrace *StackTrace     `json:"stack,omitempty"`
	Fields     ValidationError `json:"fields,omitempty"`
}

type HTTPError struct {
	err    error
	status int
	stack  []byte
}

type StackTrace struct {
	GoRoutine string        `json:"go_routine"`
	Stack     []*StackFrame `json:"stack"`
}

type StackFrame struct {
	Call  string `json:"call"`
	File  string `json:"file"`
	Line  int    `json:"line"`
	Extra int    `json:"-"`
}

var _ error = &HTTPError{}
var _ Responder = &HTTPError{}

func NewDefaultHTTPError(status int) *HTTPError {
	return &HTTPError{
		err:    errors.New(http.StatusText(status)),
		status: status,
	}
}
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

func (e *HTTPError) AddStack() {
	e.stack = debug.Stack()
}

func (e *HTTPError) Respond(w http.ResponseWriter, r *http.Request) error {

	response := errResponse{
		Error:      template.HTML(e.err.Error()),
		Status:     e.status,
		StatusText: http.StatusText(e.status),
	}
	if validationErr, ok := e.err.(ValidationError); ok {
		response.Fields = validationErr
	}

	if e.status == 500 {
		response.StackTrace = parseStack(e.stack)
	}

	if strings.HasPrefix(r.Header.Get("Accept"), "text/html") {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		t := template.New("error").Funcs(template.FuncMap{
			"jsonEncode": func(v any) string {
				b, err := json.Marshal(v)
				if err != nil {
					return err.Error()
				}
				return string(b)
			},
			"isLocal": func(s string) bool {
				return strings.HasPrefix(s, cwd)
			},
		})
		t, err = t.Parse(errorTemplate)
		if err != nil {
			panic(err)
		}

		if err, ok := e.err.(HTMLError); ok {
			response.Error = template.HTML(err.HTMLError())
		} else {
			response.Error = template.HTML("<h2>" + e.err.Error() + "</h2>")
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
		return NewResponse(reader).SetStatus(e.status).AddHeader("Content-Type", "text/html").Respond(w, r)
	}
	return NewJSONResponse(response).SetStatus(e.status).Respond(w, r)
}

func parseStack(stack []byte) *StackTrace {
	lines := bytes.Split(stack, []byte("\n"))
	goRoutine := string(lines[0])

	frames := []*StackFrame{}

	for i := 1; i < len(lines); i += 2 {
		if len(lines) <= i+1 {
			break
		}
		call := string(lines[i])
		pathLine := strings.SplitN(string(lines[i+1]), ":", 2)
		file := strings.TrimSpace(pathLine[0])
		pathLineEnd := strings.SplitN(pathLine[1], " ", 2)

		line, err := strconv.Atoi(pathLineEnd[0])
		if err != nil {
			line = -1
		}
		extra := -1
		if len(pathLineEnd) > 1 {
			c, err := parsInt(pathLineEnd[1])
			if err == nil {
				extra = c
			}
		}
		frames = append(frames, &StackFrame{
			Call:  call,
			File:  file,
			Line:  line,
			Extra: extra,
		})
	}
	return &StackTrace{
		GoRoutine: goRoutine,
		Stack:     frames,
	}
}

func parsInt(s string) (int, error) {
	sign := ""
	base := 10

	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		s = s[1:]
		sign = "-"
	}

	s, ok := strings.CutPrefix(s, "0x")
	if ok {
		base = 16
	}

	i, err := strconv.ParseInt(sign+s, base, strconv.IntSize)
	return int(i), err
}
