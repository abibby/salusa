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
	case dialects.DataTypeString, dialects.DataTypeJSON, dialects.DataTypeDate, dialects.DataTypeDateTime:
		return "VARCHAR(255)"
	case dialects.DataTypeInt32, dialects.DataTypeBoolean:
		return "INTEGER"
	case dialects.DataTypeUInt32:
		return "INTEGER UNSIGNED"
	case dialects.DataTypeFloat32:
		return "FLOAT"
	}

	return string(t)
}

func (*MySQL) CurrentTime() string {
	return "CURRENT_TIMESTAMP"
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
