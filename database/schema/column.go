package schema

import (
	"fmt"

	"gosalusa.com/database/dialects"
)

type ColumnBuilder struct {
	name     string
	datatype dialects.DataType

	nullable           bool
	primary            bool
	autoIncrement      bool
	change             bool
	defaultValue       any
	afterColumn        string
	unique             bool
	defaultCurrentTime bool
	index              bool
}

func NewColumn(name string, datatype dialects.DataType) *ColumnBuilder {
	return &ColumnBuilder{
		name:     name,
		datatype: datatype,
	}
}

func (b *ColumnBuilder) Equals(newB *ColumnBuilder) bool {
	return b.datatype == newB.datatype &&
		b.nullable == newB.nullable &&
		b.autoIncrement == newB.autoIncrement &&
		b.index == newB.index &&
		b.name == newB.name &&
		b.primary == newB.primary
}

func (b *ColumnBuilder) Name() string {
	return b.name
}

func (b *ColumnBuilder) Nullable() *ColumnBuilder {
	b.nullable = true
	return b
}
func (b *ColumnBuilder) NotNullable() *ColumnBuilder {
	b.nullable = false
	return b
}
func (b *ColumnBuilder) Primary() *ColumnBuilder {
	b.primary = true
	return b
}
func (b *ColumnBuilder) AutoIncrement() *ColumnBuilder {
	b.autoIncrement = true
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
	b.defaultValue = v
	return b
}
func (b *ColumnBuilder) Type(datatype dialects.DataType) *ColumnBuilder {
	b.datatype = datatype
	return b
}
func (b *ColumnBuilder) Unique() *ColumnBuilder {
	b.unique = true
	return b
}
func (b *ColumnBuilder) DefaultCurrentTime() *ColumnBuilder {
	b.defaultCurrentTime = true
	return b
}
func (b *ColumnBuilder) Index() *ColumnBuilder {
	b.index = true
	return b
}
func (b *ColumnBuilder) Size(s int) *ColumnBuilder {
	b.datatype.Size = s
	return b
}
func (b *ColumnBuilder) ColumnDefinition() *dialects.ColumnDefinition {
	return &dialects.ColumnDefinition{
		Name:               b.name,
		Datatype:           b.datatype,
		Nullable:           b.nullable,
		Primary:            b.primary,
		AutoIncrement:      b.autoIncrement,
		DefaultValue:       b.defaultValue,
		Unique:             b.unique,
		DefaultCurrentTime: b.defaultCurrentTime,
	}
}

func (b *ColumnBuilder) GoString() string {
	src := ""
	if b.primary {
		src += ".Primary()"
	}
	if b.autoIncrement {
		src += ".AutoIncrement()"
	}
	if b.nullable {
		src += ".Nullable()"
	}
	if b.defaultValue != nil {
		src += fmt.Sprintf(".Default(%#v)", b.defaultValue)
	}
	if b.defaultCurrentTime {
		src += ".DefaultCurrentTime()"
	}
	if b.unique {
		src += ".Unique()"
	}
	if b.index {
		src += ".Index()"
	}
	if b.change {
		src += ".Change()"
	}
	if b.datatype.Size != 0 {
		src += fmt.Sprintf(".Size(%#v)", b.datatype.Size)
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
