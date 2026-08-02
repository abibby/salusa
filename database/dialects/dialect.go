package dialects

type Dialect interface {
	EncodeSelectQuery(q *SelectQuery) (RawQuery, error)
	EncodeInsertQuery(q *InsertQuery) (RawQuery, error)
	EncodeUpdateQuery(q *UpdateQuery) (RawQuery, error)
	EncodeDeleteQuery(q *DeleteQuery) (RawQuery, error)
	EncodeCreateTableQuery(q *CreateTableQuery) (RawQuery, error)
	EncodeDropTableQuery(q *DropTableQuery) (RawQuery, error)
	EncodeAlterTableQuery(q *AlterTableQuery) (RawQuery, error)
}

func SetDefaultDialect(dialectFactory func() Dialect) {
	defaultDialect = dialectFactory
}

func New() Dialect {
	return defaultDialect()
}

var defaultDialect func() Dialect = func() Dialect { panic("no default dialect set") }
