package dialects

type Dialect interface {
	EncodeSelectQuery(q *SelectQuery) (SQLResult, error)
	EncodeUpdateQuery(q *UopdateQuery) (SQLResult, error)
}

func SetDefaultDialect(dialectFactory func() Dialect) {
	defaultDialect = dialectFactory
}

func New() Dialect {
	return defaultDialect()
}

var defaultDialect func() Dialect = func() Dialect { panic("no default dialect set") }
