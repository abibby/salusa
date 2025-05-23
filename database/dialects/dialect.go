package dialects

import "github.com/abibby/salusa/extra/sets"

type DataType string

const (
	DataTypeBlob   = DataType("blob")
	DataTypeString = DataType("string")
	DataTypeText   = DataType("text")
	DataTypeEnum   = DataType("enum")

	DataTypeBoolean = DataType("bool")

	DataTypeDate     = DataType("date")
	DataTypeDateTime = DataType("date-time")

	DataTypeFloat32 = DataType("float32")
	DataTypeFloat64 = DataType("float64")

	DataTypeInt8  = DataType("int8")
	DataTypeInt16 = DataType("int16")
	DataTypeInt32 = DataType("int32")
	DataTypeInt64 = DataType("int64")

	DataTypeUInt8  = DataType("uint8")
	DataTypeUInt16 = DataType("uint16")
	DataTypeUInt32 = DataType("uint32")
	DataTypeUInt64 = DataType("uint64")

	DataTypeJSON = DataType("json")
)

var dataTypes = sets.New(
	DataTypeBlob,
	DataTypeString,
	DataTypeText,
	DataTypeEnum,

	DataTypeBoolean,

	DataTypeDate,
	DataTypeDateTime,

	DataTypeFloat32,
	DataTypeFloat64,

	DataTypeInt8,
	DataTypeInt16,
	DataTypeInt32,
	DataTypeInt64,

	DataTypeUInt8,
	DataTypeUInt16,
	DataTypeUInt32,
	DataTypeUInt64,

	DataTypeJSON,
)

func (d DataType) IsValid() bool {
	return dataTypes.Has(d)
}

// DataTyper must not be implemented on an interface
type DataTyper interface {
	DataType() DataType
}

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
