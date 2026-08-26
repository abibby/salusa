package {{ .Package }}

import (
	"context"

	"abibby.com/salusa/database/builder"
	"abibby.com/salusa/database/model"
)

//go:generate spice generate:migration
type {{ .Name }} struct {
	model.BaseModel

	ID int `json:"id" db:"id,primary,autoincrement"`
}

func init() {
	providers.Add(modeldi.Register[*{{ .Name }}])
}

func {{ .Name }}Query(ctx context.Context) *builder.ModelBuilder[*{{ .Name }}] {
	return builder.From[*{{ .Name }}]().WithContext(ctx)
}
