package migrate

import (
	"fmt"

	"github.com/abibby/salusa/database/models"
	"github.com/abibby/salusa/internal/helpers"
)

type dropTable string

func drop(table models.Model) dropTable {
	return dropTable(helpers.GetTable(table))
}

func (dt dropTable) ToGo() string {
	return fmt.Sprintf("schema.DropIfExists(%#v)", dt)
}
