package schema

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type CreateTableBuilder struct {
	blueprint   *Blueprint
	ifNotExists bool
}

var _ helpers.SQLStringer = &CreateTableBuilder{}
var _ Blueprinter = &CreateTableBuilder{}

func Create(name string, cb func(b *Blueprint)) *CreateTableBuilder {
	b := NewBlueprint(name)
	cb(b)
	return &CreateTableBuilder{
		blueprint: b,
	}
}

func (b *CreateTableBuilder) GetBlueprint() *Blueprint {
	return b.blueprint
}
func (b *CreateTableBuilder) Type() BlueprintType {
	return BlueprintTypeCreate
}
func (b *CreateTableBuilder) SQLString(d dialects.Dialect) (string, []any, error) {
	r := helpers.Result()
	r.AddString("CREATE TABLE")
	if b.ifNotExists {
		r.AddString("IF NOT EXISTS")
	}
	r.Add(helpers.Identifier(b.blueprint.name))

	// primaryKeyColumns := slices.Filter(b.blueprint.columns, func(column *ColumnBuilder) bool {
	// 	return column.primary
	// })
	// primaryKeys := slices.Map(primaryKeyColumns, func(column *ColumnBuilder) string {
	// 	return column.name
	// })
	columns := make([]helpers.SQLStringer, len(b.blueprint.columns))
	for i, c := range b.blueprint.columns {
		columns[i] = c
	}
	if len(b.blueprint.primaryKeys) > 0 {
		columns = append(columns, helpers.Concat(
			helpers.Raw("PRIMARY KEY "),
			helpers.Group(
				helpers.Join(helpers.IdentifierList(b.blueprint.primaryKeys), ", "),
			),
		))
	}
	for _, foreignKey := range b.blueprint.foreignKeys {
		columns = append(columns, foreignKey)
		// r.Add(builder.Concat(foreignKey, builder.Raw(";")))
	}
	r.Add(helpers.Concat(
		helpers.Group(
			helpers.Concat(
				helpers.Join(columns, ", "),
			),
		),
		helpers.Raw(";"),
	))

	for _, index := range b.blueprint.indexes {
		r.Add(helpers.Concat(index, helpers.Raw(";")))
	}
	return r.SQLString(d)
}

func (b *CreateTableBuilder) GoString() string {
	return fmt.Sprintf(
		"schema.Create(%#v, %#v)",
		b.blueprint.name,
		b.blueprint,
	)
}
func (b *CreateTableBuilder) Run(ctx context.Context, tx database.DB) error {
	return runQuery(ctx, tx, b)
}
func (b *CreateTableBuilder) IfNotExists() *CreateTableBuilder {
	b.ifNotExists = true
	return b
}
