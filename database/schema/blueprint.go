package schema

import (
	"fmt"
	"strings"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/set"
	"github.com/abibby/salusa/slices"
)

type BlueprintType int

const (
	BlueprintTypeCreate BlueprintType = iota
	BlueprintTypeUpdate
)

type Blueprinter interface {
	GetBlueprint() *Blueprint
	Type() BlueprintType
}

type Blueprint struct {
	name        string
	columns     []*ColumnBuilder
	dropColumns []string
	indexes     []*IndexBuilder
	foreignKeys []*ForeignKeyBuilder
	primaryKeys []string
}

func NewBlueprint(name string) *Blueprint {
	return &Blueprint{
		name:        name,
		columns:     []*ColumnBuilder{},
		dropColumns: []string{},
		indexes:     []*IndexBuilder{},
		foreignKeys: []*ForeignKeyBuilder{},
	}
}

func (b *Blueprint) findColumn(name string) (*ColumnBuilder, bool) {
	return slices.Find(b.columns, func(c *ColumnBuilder) bool {
		return c.name == name
	})
}

func (t *Blueprint) GetBlueprint() *Blueprint {
	return t
}
func (t *Blueprint) TableName() string {
	return t.name
}

func (t *Blueprint) OfType(datatype dialects.DataType, name string) *ColumnBuilder {
	c := NewColumn(name, datatype)
	t.AddColumn(c)
	return c
}
func (t *Blueprint) AddColumn(c *ColumnBuilder) *Blueprint {
	t.columns = append(t.columns, c)
	return t
}
func (t *Blueprint) String(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeString, name)
}
func (t *Blueprint) Text(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeText, name)
}

func (t *Blueprint) Bool(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeBoolean, name)
}

func (t *Blueprint) Int(name string) *ColumnBuilder {
	return t.Int32(name)
}

func (t *Blueprint) Int8(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeInt8, name)
}

func (t *Blueprint) Int16(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeInt16, name)
}

func (t *Blueprint) Int32(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeInt32, name)
}

func (t *Blueprint) Int64(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeInt64, name)
}
func (t *Blueprint) UInt(name string) *ColumnBuilder {
	return t.UInt32(name)
}

func (t *Blueprint) UInt8(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeUInt8, name)
}

func (t *Blueprint) UInt16(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeUInt16, name)
}

func (t *Blueprint) UInt32(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeUInt32, name)
}

func (t *Blueprint) UInt64(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeUInt64, name)
}

func (t *Blueprint) Float(name string) *ColumnBuilder {
	return t.Float32(name)
}
func (t *Blueprint) Float32(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeFloat32, name)
}
func (t *Blueprint) Float64(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeFloat64, name)
}

func (t *Blueprint) JSON(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeJSON, name)
}
func (t *Blueprint) Date(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeDate, name)
}
func (t *Blueprint) DateTime(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeDateTime, name)
}
func (t *Blueprint) Blob(name string) *ColumnBuilder {
	return t.OfType(dialects.DataTypeBlob, name)
}

func (t *Blueprint) Index(name string) *IndexBuilder {
	c := newIndexBuilder(t.TableName())
	c.name = name
	t.indexes = append(t.indexes, c)
	return c
}

func (t *Blueprint) ForeignKey(localKey, relatedTable, relatedKey string) {
	f := &ForeignKeyBuilder{
		localKey:     localKey,
		relatedTable: relatedTable,
		relatedKey:   relatedKey,
	}
	t.foreignKeys = append(t.foreignKeys, f)
}

func (t *Blueprint) PrimaryKey(columns ...string) {
	t.primaryKeys = columns
}

func (t *Blueprint) DropColumn(column string) {
	t.dropColumns = append(t.dropColumns, column)
}

func (b *Blueprint) ToGo() string {
	src := "func(table *schema.Blueprint) {\n"
	for _, c := range b.columns {
		m := map[dialects.DataType]string{
			dialects.DataTypeBlob:     "Blob",
			dialects.DataTypeBoolean:  "Bool",
			dialects.DataTypeDate:     "Date",
			dialects.DataTypeDateTime: "DateTime",
			dialects.DataTypeFloat32:  "Float",
			dialects.DataTypeFloat64:  "Float64",
			dialects.DataTypeInt8:     "Int8",
			dialects.DataTypeInt16:    "Int16",
			dialects.DataTypeInt32:    "Int",
			dialects.DataTypeInt64:    "Int64",
			dialects.DataTypeJSON:     "JSON",
			dialects.DataTypeString:   "String",
			dialects.DataTypeUInt8:    "UInt8",
			dialects.DataTypeUInt16:   "UInt16",
			dialects.DataTypeUInt32:   "UInt",
			dialects.DataTypeUInt64:   "UInt64",
		}
		src += fmt.Sprintf("\ttable.%s(%#v)%s\n", m[c.datatype], c.name, c.ToGo())
	}

	for _, index := range b.indexes {
		src += fmt.Sprintf("\ttable.Index(%#v)%s\n", index.name, index.ToGo())
	}

	for _, c := range b.dropColumns {
		src += fmt.Sprintf("\ttable.DropColumn(%#v)\n", c)
	}

	for _, foreignKey := range b.foreignKeys {
		src += fmt.Sprintf("\ttable.ForeignKey(%#v, %#v, %#v)\n", foreignKey.localKey, foreignKey.relatedTable, foreignKey.relatedKey)
	}

	if len(b.primaryKeys) > 1 {
		args := strings.Join(
			slices.Map(b.primaryKeys, func(pKey string) string {
				return fmt.Sprintf("%#v", pKey)
			}),
			", ",
		)
		src += fmt.Sprintf("\ttable.PrimaryKey(%s)\n", args)
	}

	return src + "}"
}

func (t *Blueprint) Merge(newBlueprint *Blueprint) {
	if t.name != newBlueprint.name {
		return
	}

	for _, newColumn := range newBlueprint.columns {
		if newColumn.change {
			for i, c := range t.columns {
				if c.name == newColumn.name {
					t.columns[i] = newColumn
					break
				}
			}
		} else {
			t.columns = append(t.columns, newColumn)
		}
	}

	t.columns = slices.Filter(t.columns, func(c *ColumnBuilder) bool {
		return !slices.Has(newBlueprint.dropColumns, c.name)
	})

	t.foreignKeys = append(t.foreignKeys, newBlueprint.foreignKeys...)
	t.indexes = append(t.indexes, newBlueprint.indexes...)
	if newBlueprint.primaryKeys != nil {
		t.primaryKeys = newBlueprint.primaryKeys
	}
}

func (t *Blueprint) Update(oldBlueprint, newBlueprint *Blueprint) bool {
	addedColumns := set.New[string]()
	hasChanges := false
	for _, newColumn := range newBlueprint.columns {
		oldColumn, ok := oldBlueprint.findColumn(newColumn.name)
		if ok {
			addedColumns.Add(newColumn.name)
			if newColumn.Equals(oldColumn) {
				continue
			}
		}

		newColumn.change = ok
		hasChanges = true
		t.AddColumn(newColumn)
	}
	for _, oldColumn := range oldBlueprint.columns {
		if !addedColumns.Has(oldColumn.name) {
			hasChanges = true
			t.DropColumn(oldColumn.name)
		}
	}

	for _, newKey := range newBlueprint.foreignKeys {
		_, ok := slices.Find(oldBlueprint.foreignKeys, func(oldKey *ForeignKeyBuilder) bool {
			return newKey.localKey == oldKey.localKey &&
				newKey.relatedKey == oldKey.relatedKey &&
				newKey.relatedTable == oldKey.relatedTable
		})
		if !ok {
			t.foreignKeys = append(t.foreignKeys, newKey)
			hasChanges = true
		}
	}
	for _, newIndex := range newBlueprint.indexes {
		_, ok := slices.Find(oldBlueprint.indexes, func(oldIndex *IndexBuilder) bool {
			return newIndex.name == oldIndex.name
		})
		if !ok {
			t.indexes = append(t.indexes, newIndex)
			hasChanges = true
		}
	}

	// TODO: add support for primary keys
	return hasChanges
}
