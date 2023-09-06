package migrate

import (
	"fmt"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/models"
	"github.com/abibby/salusa/database/schema"
	"github.com/abibby/salusa/internal/helpers"
)

var (
	ErrNoChanges = fmt.Errorf("no changes")
)

func (m *Migrations) update(model models.Model) (*schema.UpdateTableBuilder, *schema.UpdateTableBuilder, error) {
	err := builder.InitializeRelationships(model)
	if err != nil {
		return nil, nil, err
	}

	tableName := helpers.GetTable(model)
	fields, err := getFields(model)
	if err != nil {
		return nil, nil, err
	}
	oldTable := m.Blueprint(tableName)
	newTable := blueprintFromFields(tableName, fields)

	hasChanges := false

	up := schema.Table(tableName, func(table *schema.Blueprint) {
		hasChanges = table.Update(oldTable, newTable)
	})
	down := schema.Table(tableName, func(table *schema.Blueprint) {
		table.Update(newTable, oldTable)
	})
	if !hasChanges {
		return nil, nil, ErrNoChanges
	}

	return up, down, nil
}
