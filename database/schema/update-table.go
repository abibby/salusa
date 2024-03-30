package schema

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type UpdateTableBuilder struct {
	blueprint *Blueprint
}

var _ helpers.SQLStringer = &UpdateTableBuilder{}
var _ Blueprinter = &UpdateTableBuilder{}

func Table(name string, cb func(table *Blueprint)) *UpdateTableBuilder {
	table := NewBlueprint(name)
	cb(table)
	return &UpdateTableBuilder{
		blueprint: table,
	}
}

func (b *UpdateTableBuilder) GetBlueprint() *Blueprint {
	return b.blueprint
}
func (b *UpdateTableBuilder) Type() BlueprintType {
	return BlueprintTypeUpdate
}

func (b *UpdateTableBuilder) SQLString(d dialects.Dialect) (string, []any, error) {
	r := helpers.Result()
	alterTable := helpers.Concat(helpers.Raw("ALTER TABLE "), helpers.Identifier(b.blueprint.name))
	for _, column := range b.blueprint.dropColumns {
		r.Add(helpers.Concat(
			alterTable,
			helpers.Raw(" DROP COLUMN "),
			helpers.Identifier(column),
			helpers.Raw(";"),
		))
	}
	for _, column := range b.blueprint.columns {
		if column.change {
			r.Add(helpers.Concat(
				alterTable,
				helpers.Raw(" MODIFY COLUMN "),
				column,
				helpers.Raw(";"),
			))
		} else {
			r.Add(helpers.Concat(
				alterTable,
				helpers.Raw(" ADD "),
				column,
				helpers.Raw(";"),
			))
		}
	}
	for _, foreignKey := range b.blueprint.foreignKeys {
		r.Add(helpers.Concat(
			alterTable,
			helpers.Raw(" ADD "),
			foreignKey,
			helpers.Raw(";"),
		))
	}
	for _, index := range b.blueprint.indexes {
		r.Add(helpers.Concat(index, helpers.Raw(";")))
	}

	return r.SQLString(d)
}

func (b *UpdateTableBuilder) GoString() string {
	return fmt.Sprintf(
		"schema.Table(%#v, %#v)",
		b.blueprint.name,
		b.blueprint,
	)
}

func (b *UpdateTableBuilder) Run(ctx context.Context, tx database.DB) error {
	return runQuery(ctx, tx, b)
}
