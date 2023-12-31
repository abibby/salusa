package mysql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/abibby/salusa/database/dialects"
)

type MySQL struct{}

func (*MySQL) Identifier(s string) string {
	parts := strings.Split(s, ".")
	for i, p := range parts {
		parts[i] = "`" + p + "`"
	}
	return strings.Join(parts, ".")
}

func (*MySQL) DataType(t dialects.DataType) string {
	switch t {
	case dialects.DataTypeString:
		return "VARCHAR(255)"
	case dialects.DataTypeText, dialects.DataTypeJSON:
		return "MEDIUMTEXT"
	case dialects.DataTypeInt8:
		return "TINYINT"
	case dialects.DataTypeInt16:
		return "SMALLINT"
	case dialects.DataTypeInt32:
		return "INT"
	case dialects.DataTypeInt64:
		return "BIGINT"
	case dialects.DataTypeUInt8:
		return "TINYINT UNSIGNED"
	case dialects.DataTypeUInt16:
		return "SMALLINT UNSIGNED"
	case dialects.DataTypeUInt32:
		return "INT UNSIGNED"
	case dialects.DataTypeUInt64:
		return "BIGINT UNSIGNED"
	case dialects.DataTypeBoolean:
		return "BOOLEAN"
	case dialects.DataTypeFloat32:
		return "FLOAT"
	case dialects.DataTypeFloat64:
		return "DOUBLE"
	case dialects.DataTypeDate:
		return "DATE"
	case dialects.DataTypeDateTime:
		return "DATETIME"
	}

	return string(t)
}

func (*MySQL) CurrentTime() string {
	return "CURRENT_TIMESTAMP"
}
func (*MySQL) AutoIncrement() string {
	return "AUTO_INCREMENT"
}

func (s *MySQL) Escape(v any) string {
	switch v := v.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprint(v)
	case bool:
		if v {
			return "1"
		} else {
			return "0"
		}
	default:
		b, err := json.Marshal(v)
		if err != nil {
			panic(fmt.Errorf("failed to escape value: %w", err))
		}
		return s.Escape(string(b))
	}
}

func (*MySQL) Binding() string {
	return "?"
}
func UseMySql() {
	dialects.SetDefaultDialect(func() dialects.Dialect {
		return &MySQL{}
	})
}
func init() {
	UseMySql()
}
