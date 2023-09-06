package schema

import (
	"fmt"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/internal/helpers"
)

type IndexBuilder struct {
	table   string
	name    string
	columns []string
	unique  bool
}

func newIndexBuilder(table string) *IndexBuilder {
	return &IndexBuilder{
		columns: []string{},
		table:   table,
	}
}

func (b *IndexBuilder) AddColumn(c string) *IndexBuilder {
	b.columns = append(b.columns, c)
	return b
}

func (b *IndexBuilder) Unique() *IndexBuilder {
	b.unique = true
	return b
}
func (b *IndexBuilder) ToSQL(d dialects.Dialect) (string, []any, error) {
	r := helpers.Result().AddString("CREATE")
	if b.unique {
		r.AddString("UNIQUE")
	}
	r.AddString("INDEX IF NOT EXISTS").
		Add(helpers.Identifier(b.name)).
		AddString("ON").
		Add(helpers.Identifier(b.table)).
		Add(helpers.Group(helpers.Join(helpers.IdentifierList(b.columns), ", ")))

	return r.ToSQL(d)
}

func (b *IndexBuilder) ToGo() string {
	src := ""
	for _, c := range b.columns {
		src += fmt.Sprintf(".AddColumn(%#v)", c)
	}
	if b.unique {
		src += ".Unique()"
	}
	return src
}
