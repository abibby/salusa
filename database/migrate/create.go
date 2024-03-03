package migrate

import (
	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/database/schema"
	"github.com/abibby/salusa/internal/relationship"
	"github.com/abibby/salusa/slices"
)

func CreateFromModel(m model.Model) (*schema.CreateTableBuilder, error) {
	err := relationship.InitializeRelationships(m)
	if err != nil {
		panic(err)
	}

	tableName := database.GetTable(m)
	fields, err := getFields(m)
	if err != nil {
		return nil, err
	}

	return schema.Create(tableName, func(table *schema.Blueprint) {
		table.Merge(blueprintFromFields(tableName, fields))
	}), nil
}

func blueprintFromFields(tableName string, fields []*field) *schema.Blueprint {
	table := schema.NewBlueprint(tableName)

	addedForeignKeys := []*builder.ForeignKey{}
	primaryColumns := []*schema.ColumnBuilder{}
	for _, f := range fields {
		if f.relation != nil {
			foreignKeys := f.relation.ForeignKeys()
			for _, foreignKey := range foreignKeys {
				if slices.Has(addedForeignKeys, foreignKey) {
					continue
				}

				addedForeignKeys = append(addedForeignKeys, foreignKey)
				table.ForeignKey(foreignKey.LocalKey, foreignKey.RelatedTable, foreignKey.RelatedKey)
			}
		} else {
			b := table.OfType(f.dataType, f.tag.Name)
			if f.nullable {
				b.Nullable()
			}
			if f.tag.Index {
				b.Index()
				// table.Index(fmt.Sprintf("%s-%s", tableName, field.tag.Name)).AddColumn(field.tag.Name)
			}
			if f.tag.AutoIncrement {
				b.AutoIncrement()
			}
			if f.tag.Primary {
				primaryColumns = append(primaryColumns, b)
			}
			if f.tag.Unique {
				b.Unique()
			}
		}
	}

	if len(primaryColumns) > 1 {
		table.PrimaryKey(slices.Map(primaryColumns, func(c *schema.ColumnBuilder) string {
			return c.Name()
		})...)
	} else {
		for _, b := range primaryColumns {
			b.Primary()
		}
	}
	return table
}
