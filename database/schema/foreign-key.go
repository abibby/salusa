package schema

import (
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type ForeignKeyBuilder struct {
	relatedTable string
	localKey     string
	relatedKey   string
}

func (b *ForeignKeyBuilder) id() string {
	return b.localKey + "-" + b.relatedTable + "-" + b.relatedKey
}

func (b *ForeignKeyBuilder) ToSQL(d dialects.Dialect) (string, []any, error) {
	r := helpers.Result()

	r.AddString("CONSTRAINT").
		Add(helpers.Identifier(b.id())).
		AddString("FOREIGN KEY").
		Add(helpers.Group(helpers.Identifier(b.localKey))).
		AddString("REFERENCES").
		Add(helpers.Concat(
			helpers.Identifier(b.relatedTable),
			helpers.Group(helpers.Identifier(b.relatedKey)),
		))
	return r.ToSQL(d)
}
