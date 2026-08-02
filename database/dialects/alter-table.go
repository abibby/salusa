package dialects

type AlterTableQueryBuilder interface {
	AlterTableQuery() *AlterTableQuery
}

type AlterTableQuery struct {
	Table         string
	DropColumns   []string
	ModifyColumns []ColumnDefinition
	AddColumns    []ColumnDefinition
	ForeignKeys   []ForeignKey
	Indexes       []Index
}
