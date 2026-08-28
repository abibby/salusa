package schema

import "abibby.com/salusa/database/dialects"

type ForeignKeyBuilder struct {
	relatedTable string
	localKey     string
	relatedKey   string
}

func (b *ForeignKeyBuilder) ForeignKey() *dialects.ForeignKey {
	return &dialects.ForeignKey{
		Name:           b.localKey + "-" + b.relatedTable + "-" + b.relatedKey,
		Columns:        []string{b.localKey},
		ForeignTable:   b.relatedTable,
		ForeignColumns: []string{b.relatedKey},
	}
}
