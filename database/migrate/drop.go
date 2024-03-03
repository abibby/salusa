package migrate

import (
	"fmt"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/model"
)

type dropTable string

func drop(table model.Model) dropTable {
	return dropTable(database.GetTable(table))
}

func (dt dropTable) ToGo() string {
	return fmt.Sprintf("schema.DropIfExists(%#v)", dt)
}
