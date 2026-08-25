package schema

import (
	"fmt"
	"slices"
	"strings"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/extra/sets"
	"github.com/abibby/salusa/stream"
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
	return stream.OfSlice(b.columns).Find(func(c *ColumnBuilder) bool {
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

func (b *Blueprint) GoString() string {
	src := strings.Builder{}
	src.WriteString("func(table *schema.Blueprint) {\n")
	for _, c := range b.columns {
		m := map[string]string{
			dialects.DataTypeBlob.Name:     "Blob",
			dialects.DataTypeBoolean.Name:  "Bool",
			dialects.DataTypeDate.Name:     "Date",
			dialects.DataTypeDateTime.Name: "DateTime",
			dialects.DataTypeFloat32.Name:  "Float",
			dialects.DataTypeFloat64.Name:  "Float64",
			dialects.DataTypeInt8.Name:     "Int8",
			dialects.DataTypeInt16.Name:    "Int16",
			dialects.DataTypeInt32.Name:    "Int",
			dialects.DataTypeInt64.Name:    "Int64",
			dialects.DataTypeJSON.Name:     "JSON",
			dialects.DataTypeString.Name:   "String",
			dialects.DataTypeUInt8.Name:    "UInt8",
			dialects.DataTypeUInt16.Name:   "UInt16",
			dialects.DataTypeUInt32.Name:   "UInt",
			dialects.DataTypeUInt64.Name:   "UInt64",
		}
		fmt.Fprintf(&src, "\ttable.%s(%#v)%s\n", m[c.datatype.Name], c.name, c.GoString())
	}

	for _, index := range b.indexes {
		fmt.Fprintf(&src, "\ttable.Index(%#v)%s\n", index.name, index.GoString())
	}

	for _, c := range b.dropColumns {
		fmt.Fprintf(&src, "\ttable.DropColumn(%#v)\n", c)
	}

	for _, foreignKey := range b.foreignKeys {
		fmt.Fprintf(&src, "\ttable.ForeignKey(%#v, %#v, %#v)\n", foreignKey.localKey, foreignKey.relatedTable, foreignKey.relatedKey)
	}

	if len(b.primaryKeys) > 1 {
		args := strings.Join(
			stream.OfSlice(b.primaryKeys).Map(func(pKey string) string {
				return fmt.Sprintf("%#v", pKey)
			}).Slice(),
			", ",
		)
		fmt.Fprintf(&src, "\ttable.PrimaryKey(%s)\n", args)
	}

	src.WriteString("}")
	return src.String()
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

	t.columns = stream.OfSlice(t.columns).Filter(func(c *ColumnBuilder) bool {
		return !slices.Contains(newBlueprint.dropColumns, c.name)
	}).Slice()

	t.foreignKeys = append(t.foreignKeys, newBlueprint.foreignKeys...)
	t.indexes = append(t.indexes, newBlueprint.indexes...)
	if newBlueprint.primaryKeys != nil {
		t.primaryKeys = newBlueprint.primaryKeys
	}
}

func (t *Blueprint) Update(oldBlueprint, newBlueprint *Blueprint) bool {
	addedColumns := sets.New[string]()
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
		_, ok := stream.OfSlice(oldBlueprint.foreignKeys).Find(func(oldKey *ForeignKeyBuilder) bool {
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
		_, ok := stream.OfSlice(oldBlueprint.indexes).Find(func(oldIndex *IndexBuilder) bool {
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
