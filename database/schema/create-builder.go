package schema

import (
	"context"
	"fmt"

	"abibby.com/salusa/database"
	"abibby.com/salusa/database/dialects"
)

type CreateTableBuilder struct {
	blueprint   *Blueprint
	ifNotExists bool
	temporary   bool
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

func (b *CreateTableBuilder) CreateTableQuery() *dialects.CreateTableQuery {
	columns := make([]dialects.ColumnDefinition, len(b.blueprint.columns))
	for i, c := range b.blueprint.columns {
		columns[i] = *c.ColumnDefinition()
	}

	foreignKeys := make([]dialects.ForeignKey, len(b.blueprint.foreignKeys))
	for i, fk := range b.blueprint.foreignKeys {
		foreignKeys[i] = *fk.ForeignKey()
	}

	indexes := make([]dialects.Index, len(b.blueprint.indexes))
	for i, fk := range b.blueprint.indexes {
		indexes[i] = *fk.Index()
	}

	return &dialects.CreateTableQuery{
		IfNotExists: b.ifNotExists,
		Temporary:   b.temporary,
		Table:       b.blueprint.TableName(),
		Columns:     columns,
		PrimaryKeys: b.blueprint.primaryKeys,
		ForeignKeys: foreignKeys,
		Indexes:     indexes,
	}
}

func (b *CreateTableBuilder) GoString() string {
	return fmt.Sprintf(
		"schema.Create(%#v, %#v)",
		b.blueprint.name,
		b.blueprint,
	)
}
func (b *CreateTableBuilder) Run(ctx context.Context, tx database.DB) error {
	q := b.CreateTableQuery()
	d, err := dialects.New(tx.DriverName())
	if err != nil {
		return err
	}
	result, err := d.EncodeCreateTableQuery(q)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, result.SQL, result.Bindings...)
	return err
}
func (b *CreateTableBuilder) IfNotExists() *CreateTableBuilder {
	b.ifNotExists = true
	return b
}
func (b *CreateTableBuilder) Temporary() *CreateTableBuilder {
	b.temporary = true
	return b
}
