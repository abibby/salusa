package schema

import (
	"fmt"
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
