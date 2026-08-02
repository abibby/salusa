package dialects

type CreateTableQueryBuilder interface {
	CreateTableQuery() *CreateTableQuery
}
type CreateTableQuery struct {
	IfNotExists bool
	Table       string
	Columns     []ColumnDefinition
}

type ColumnDefinition struct {
	Name               string
	Datatype           DataType
	Nullable           bool
	Primary            bool
	AutoIncrement      bool
	DefaultValue       any
	Unique             bool
	DefaultCurrentTime bool
}
