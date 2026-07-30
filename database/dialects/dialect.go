package dialects

type Dialect interface {
	EncodeQuery(q *Query) (SQLResult, error)
}

func SetDefaultDialect(dialectFactory func() Dialect) {
	defaultDialect = dialectFactory
}

func New() Dialect {
	return defaultDialect()
}

var defaultDialect func() Dialect = func() Dialect { panic("no default dialect set") }
