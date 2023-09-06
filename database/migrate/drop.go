package migrate

import (
	"fmt"

	"github.com/abibby/salusa/database/internal/helpers"
	"github.com/abibby/salusa/database/models"
)

type dropTable string

func drop(table models.Model) dropTable {
	return dropTable(helpers.GetTable(table))
}

func (dt dropTable) ToGo() string {
	return fmt.Sprintf("schema.DropIfExists(%#v)", dt)
}
