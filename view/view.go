package view

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"io/fs"
	"net/http"

	"github.com/abibby/salusa/di"
	"github.com/abibby/salusa/request"
	"github.com/abibby/salusa/router"
)

type ViewTemplate struct {
	fsys     fs.FS
	patterns []string
}

type ViewRequest struct {
	URL      router.URLResolver `inject:""`
	Template *ViewTemplate      `inject:""`
}
type ViewResponse struct {
	req  *ViewRequest
	file string
	data any
}

var _ request.Responder = (*ViewResponse)(nil)

func (vr *ViewResponse) Execute(w io.Writer) error {
	tpl, err := template.New("").
		Funcs(template.FuncMap{
			"route": vr.req.URL.Resolve,
		}).
		ParseFS(vr.req.Template.fsys, vr.req.Template.patterns...)
	if err != nil {
		return err
	}
	return tpl.ExecuteTemplate(w, vr.file, vr.data)
}

func (vr *ViewResponse) Respond(w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Content-Type", "text/html")
	return vr.Execute(w)
}

func (vr *ViewResponse) Bytes() ([]byte, error) {
	b := &bytes.Buffer{}
	err := vr.Execute(b)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func View(file string, data any) *request.RequestHandler[ViewRequest, *ViewResponse] {
	return request.Handler(func(r *ViewRequest) (*ViewResponse, error) {
		return &ViewResponse{
			req:  r,
			file: file,
			data: data,
		}, nil
	})
}
func NewViewTemplate(fsys fs.FS, patterns ...string) *ViewTemplate {
	if len(patterns) == 0 {
		patterns = []string{"**/*.html"}
	}
	return &ViewTemplate{
		fsys:     fsys,
		patterns: patterns,
	}
}
func Register(fsys fs.FS, patterns ...string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		di.RegisterSingleton(ctx, func() *ViewTemplate {
			return NewViewTemplate(fsys, patterns...)
		})
		return nil
	}
}
