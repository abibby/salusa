package dialects

type CreateTableQueryBuilder interface {
	CreateTableQuery() *CreateTableQuery
}
type CreateTableQuery struct {
	IfNotExists bool
	Temporary   bool
	Table       string
	Columns     []ColumnDefinition
	PrimaryKeys []string
	ForeignKeys []ForeignKey
	Indexes     []Index
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

type ForeignKey struct {
	Name           string
	Columns        []string
	ForeignTable   string
	ForeignColumns []string
	// OnUpdate       string
	// OnDelete       string
}
