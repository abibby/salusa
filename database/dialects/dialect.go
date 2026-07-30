package dialects

type Dialect interface {
	Identifier(string) string
	DataType(DataType) string
	CurrentTime() string
	AutoIncrement() string
	Escape(v any) string
	Binding() string
}

type unsetDialect struct{}

func (*unsetDialect) Identifier(s string) string {
	return s
}

func (*unsetDialect) DataType(t DataType) string {
	return string(t)
}

func (*unsetDialect) CurrentTime() string {
	return "CURRENT_TIMESTAMP"
}

func (*unsetDialect) AutoIncrement() string {
	return "AUTO_INCREMENT"
}

func (*unsetDialect) Escape(v any) string {
	return ""
}

func (*unsetDialect) Binding() string {
	return "?"
}

func SetDefaultDialect(dialectFactory func() Dialect) {
	defaultDialect = dialectFactory
}

func New() Dialect {
	return defaultDialect()
}

var defaultDialect func() Dialect = func() Dialect { return &unsetDialect{} }
