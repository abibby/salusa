package dialects

type AlterTableQueryBuilder interface {
	AlterTableQuery() *AlterTableQuery
}

type AlterTableQuery struct{}
