package schema

import (
	"fmt"

	"github.com/abibby/salusa/database/dialects"
)

type ColumnBuilder struct {
	def         dialects.ColumnDefinition
	index       bool
	change      bool
	afterColumn string
}

func NewColumn(name string, datatype dialects.DataType) *ColumnBuilder {
	return &ColumnBuilder{
		def: dialects.ColumnDefinition{
			Name:     name,
			Datatype: datatype,
		},
	}
}

func (b *ColumnBuilder) Equals(newB *ColumnBuilder) bool {
	return b.def.Datatype == newB.def.Datatype &&
		b.def.Nullable == newB.def.Nullable &&
		b.def.AutoIncrement == newB.def.AutoIncrement &&
		b.index == newB.index &&
		b.def.Name == newB.def.Name &&
		b.def.Primary == newB.def.Primary
}

func (b *ColumnBuilder) Name() string {
	return b.def.Name
}

func (b *ColumnBuilder) Nullable() *ColumnBuilder {
	b.def.Nullable = true
	return b
}
func (b *ColumnBuilder) NotNullable() *ColumnBuilder {
	b.def.Nullable = false
	return b
}
func (b *ColumnBuilder) Primary() *ColumnBuilder {
	b.def.Primary = true
	return b
}
func (b *ColumnBuilder) AutoIncrement() *ColumnBuilder {
	b.def.AutoIncrement = true
	return b
}
func (b *ColumnBuilder) After(column string) *ColumnBuilder {
	b.afterColumn = column
	return b
}
func (b *ColumnBuilder) Change() *ColumnBuilder {
	b.change = true
	return b
}
func (b *ColumnBuilder) Default(v any) *ColumnBuilder {
	b.def.DefaultValue = v
	return b
}
func (b *ColumnBuilder) Type(datatype dialects.DataType) *ColumnBuilder {
	b.def.Datatype = datatype
	return b
}
func (b *ColumnBuilder) Unique() *ColumnBuilder {
	b.def.Unique = true
	return b
}
func (b *ColumnBuilder) DefaultCurrentTime() *ColumnBuilder {
	b.def.DefaultCurrentTime = true
	return b
}
func (b *ColumnBuilder) Index() *ColumnBuilder {
	b.index = true
	return b
}

func (b *ColumnBuilder) GoString() string {
	src := ""
	if b.def.Primary {
		src += ".Primary()"
	}
	if b.def.AutoIncrement {
		src += ".AutoIncrement()"
	}
	if b.def.Nullable {
		src += ".Nullable()"
	}
	if b.def.DefaultValue != nil {
		src += fmt.Sprintf(".Default(%#v)", b.def.DefaultValue)
	}
	if b.def.DefaultCurrentTime {
		src += ".DefaultCurrentTime()"
	}
	if b.def.Unique {
		src += ".Unique()"
	}
	if b.index {
		src += ".Index()"
	}
	if b.change {
		src += ".Change()"
	}
	return src
}

func (b *CreateTableBuilder) Columns(columns ...*ColumnBuilder) *CreateTableBuilder {
	b.blueprint.columns = columns
	return b
}
func (b *CreateTableBuilder) AddColumns(columns ...*ColumnBuilder) *CreateTableBuilder {
	b.blueprint.columns = append(b.blueprint.columns, columns...)
	return b
}
