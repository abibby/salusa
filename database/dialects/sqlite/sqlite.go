package sqlite

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/abibby/salusa/database/dialects"
)

type SQLite struct{}

func (*SQLite) Identifier(s string) string {
	parts := strings.Split(s, ".")
	for i, p := range parts {
		parts[i] = "\"" + p + "\""
	}
	return strings.Join(parts, ".")
}

func (*SQLite) DataType(t dialects.DataType) string {
	switch t {
	case dialects.DataTypeString, dialects.DataTypeText, dialects.DataTypeJSON:
		return "TEXT"
	case dialects.DataTypeDate, dialects.DataTypeDateTime:
		return "TIMESTAMP"
	case dialects.DataTypeInt32, dialects.DataTypeUInt32, dialects.DataTypeBoolean:
		return "INTEGER"
	case dialects.DataTypeFloat32:
		return "FLOAT"
	}
	return string(t)
}

func (*SQLite) CurrentTime() string {
	return "CURRENT_TIMESTAMP"
}

func (*SQLite) AutoIncrement() string {
	return "AUTOINCREMENT"
}

func (s *SQLite) Escape(v any) string {
	if marshaler, ok := v.(encoding.TextMarshaler); ok {
		str, err := marshaler.MarshalText()
		if err != nil {
			panic(fmt.Errorf("failed to escape value: %w", err))
		}
		return s.Escape(string(str))
	}
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.String {
		return "'" + strings.ReplaceAll(val.String(), "'", "''") + "'"
	}
	if val.CanInt() || val.CanUint() || val.CanFloat() {
		return fmt.Sprint(v)
	}
	if val.Kind() == reflect.Bool {
		if val.Bool() {
			return "1"
		} else {
			return "0"
		}
	}

	b, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Errorf("failed to escape value: %w", err))
	}
	return s.Escape(string(b))
}

func (*SQLite) Binding() string {
	return "?"
}

func UseSQLite() {
	dialects.SetDefaultDialect(func() dialects.Dialect {
		return &SQLite{}
	})
}
func init() {
	UseSQLite()
}
