package schema

import (
	"fmt"

	"github.com/abibby/salusa/database/dialects"
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

func (b *IndexBuilder) Index() *dialects.Index {
	return &dialects.Index{
		Table:   b.table,
		Name:    b.name,
		Columns: b.columns,
		Unique:  b.unique,
	}
}

func (b *IndexBuilder) GoString() string {
	src := ""
	for _, c := range b.columns {
		src += fmt.Sprintf(".AddColumn(%#v)", c)
	}
	if b.unique {
		src += ".Unique()"
	}
	return src
}
