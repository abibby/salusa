package schema

import (
	"context"
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/dialects"
)

type CreateTableBuilder struct {
	blueprint   *Blueprint
	ifNotExists bool
}

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

func (b *CreateTableBuilder) GoString() string {
	return fmt.Sprintf(
		"schema.Create(%#v, %#v)",
		b.blueprint.name,
		b.blueprint,
	)
}
func (b *CreateTableBuilder) Run(ctx context.Context, tx database.DB) error {
	q := b.blueprint.CreateTableQuery()
	result, err := dialects.New().EncodeCreateTableQuery(q)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, result.Query, result.Bindings...)
	return err
}
func (b *CreateTableBuilder) IfNotExists() *CreateTableBuilder {
	b.ifNotExists = true
	return b
}
