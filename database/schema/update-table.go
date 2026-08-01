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
	return &dialects.AlterTableQuery{}
}
func (b *UpdateTableBuilder) GoString() string {
	return fmt.Sprintf(
		"schema.Table(%#v, %#v)",
		b.blueprint.name,
		b.blueprint,
	)
}

func (b *UpdateTableBuilder) Run(ctx context.Context, tx database.DB) error {
	result, err := dialects.New().EncodeAlterTableQuery(b.AlterTableQuery())
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, result.Query, result.Bindings...)
	return err
}
