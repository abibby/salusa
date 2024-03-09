package {{ .Package }}

import (
	"context"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/database/model/modeldi"
)

//go:generate spice generate:migration
type {{ .Name }} struct {
	model.BaseModel

	ID int `json:"id" db:"id,primary,autoincrement"`
}

func {{ .Name }}Query(ctx context.Context) *builder.Builder[*{{ .Name }}] {
	return builder.From[*{{ .Name }}]().WithContext(ctx)
}
