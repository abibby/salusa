package schema

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
)

type UpdateTableBuilder struct {
	blueprint *Blueprint
}

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

func (b *UpdateTableBuilder) AlterTableQuery() *dialects.AlterTableQuery {
	addColumns := make([]dialects.ColumnDefinition, 0, len(b.blueprint.columns))
	modifyColumns := make([]dialects.ColumnDefinition, 0, len(b.blueprint.columns))

	for _, column := range b.blueprint.columns {
		if column.change {
			modifyColumns = append(modifyColumns, *column.ColumnDefinition())
		} else {
			addColumns = append(addColumns, *column.ColumnDefinition())
		}
	}
	foreignKeys := make([]dialects.ForeignKey, 0, len(b.blueprint.foreignKeys))
	for _, foreignKey := range b.blueprint.foreignKeys {
		foreignKeys = append(foreignKeys, *foreignKey.ForeignKey())
	}
	indexes := make([]dialects.Index, 0, len(b.blueprint.indexes))
	for _, index := range b.blueprint.indexes {
		indexes = append(indexes, *index.Index())
	}

	return &dialects.AlterTableQuery{
		Table:         b.blueprint.TableName(),
		DropColumns:   b.blueprint.dropColumns,
		AddColumns:    addColumns,
		ModifyColumns: modifyColumns,
		ForeignKeys:   foreignKeys,
		Indexes:       indexes,
	}
}
func (b *UpdateTableBuilder) GoString() string {
	return fmt.Sprintf(
		"schema.Table(%#v, %#v)",
		b.blueprint.name,
		b.blueprint,
	)
}

func (b *UpdateTableBuilder) Run(ctx context.Context, tx database.DB) error {
	d, err := dialects.New(tx.DriverName())
	if err != nil {
		return err
	}
	result, err := d.EncodeAlterTableQuery(b.AlterTableQuery())
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, result.SQL, result.Bindings...)
	return err
}
