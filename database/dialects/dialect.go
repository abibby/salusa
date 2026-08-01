package dialects

type Dialect interface {
	EncodeSelectQuery(q *SelectQuery) (SQLResult, error)
	EncodeInsertQuery(q *InsertQuery) (SQLResult, error)
	EncodeUpdateQuery(q *UpdateQuery) (SQLResult, error)
	EncodeDeleteQuery(q *DeleteQuery) (SQLResult, error)
	EncodeCreateTableQuery(q *CreateTableQuery) (SQLResult, error)
	EncodeDropTableQuery(q *DropTableQuery) (SQLResult, error)
	EncodeAlterTableQuery(q *AlterTableQuery) (SQLResult, error)
}

func SetDefaultDialect(dialectFactory func() Dialect) {
	defaultDialect = dialectFactory
}

func New() Dialect {
	return defaultDialect()
}

var defaultDialect func() Dialect = func() Dialect { panic("no default dialect set") }
