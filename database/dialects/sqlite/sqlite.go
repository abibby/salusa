package sqlite

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/abibby/salusa/database/dialects"
	"github.com/abibby/salusa/database/dialects/generic"
)

type SQLiteCore struct{}

func New() dialects.Dialect {
	return generic.New(&SQLiteCore{})
}

func (*SQLiteCore) Identifier(s string) string {
	if s == "*" {
		return s
	}
	parts := strings.Split(s, ".")
	for i, p := range parts {
		if p == "*" {
			continue
		}
		parts[i] = `"` + p + `"`
	}
	return strings.Join(parts, ".")
}

func (*SQLiteCore) DataType(t dialects.DataType) string {
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

func (*SQLiteCore) CurrentTime() string {
	return "CURRENT_TIMESTAMP"
}

func (*SQLiteCore) AutoIncrement() string {
	return "PRIMARY KEY AUTOINCREMENT"
}

func (s *SQLiteCore) Escape(v any) string {
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

func (*SQLiteCore) Binding() string {
	return "?"
}

func (s *SQLiteCore) Features() dialects.Features {
	return dialects.Features{
		Returning: true,
	}
}

func UseSQLite() {
	dialects.Register("sqlite3", New)
	dialects.Register("sqlite", New)
}
func init() {
	UseSQLite()
}
